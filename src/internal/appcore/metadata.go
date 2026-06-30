package appcore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const embeddedMetadataKeyword = "ClipForVRChat:AutoCapture"

type AutoCaptureEmbeddedMetadata struct {
	SchemaVersion      int                        `json:"schema_version"`
	App                string                     `json:"app"`
	BatchID            string                     `json:"batch_id"`
	ShotID             string                     `json:"shot_id"`
	CapturedAtUTC      string                     `json:"captured_at_utc"`
	CaptureMode        string                     `json:"capture_mode"`
	View               CameraViewConfig           `json:"view"`
	Stream             *AutoCaptureStreamMetadata `json:"stream,omitempty"`
	PresenceSource     string                     `json:"presence_source"`
	PresenceConfidence string                     `json:"presence_confidence"`
	Users              []PresenceUser             `json:"users,omitempty"`
}

func BuildAutoCaptureEmbeddedMetadata(cfg AutoCaptureConfig, batchID string, shotID string, view CameraViewConfig, users []PresenceUser, confidence string, streamInfo SpoutCaptureResult) AutoCaptureEmbeddedMetadata {
	metadataUsers := []PresenceUser(nil)
	if cfg.Output.WriteUserListToEXIF {
		metadataUsers = users
		if !cfg.Output.WriteUserIDsToEXIF {
			metadataUsers = presenceUsersWithoutIDs(users)
		}
	}
	return AutoCaptureEmbeddedMetadata{
		SchemaVersion:      1,
		App:                "ClipForVRChat",
		BatchID:            batchID,
		ShotID:             shotID,
		CapturedAtUTC:      time.Now().UTC().Format(time.RFC3339),
		CaptureMode:        cfg.Capture.Mode,
		View:               view,
		Stream:             autoCaptureStreamMetadata(streamInfo),
		PresenceSource:     "output_log",
		PresenceConfidence: confidence,
		Users:              metadataUsers,
	}
}

func WriteAutoCaptureEmbeddedMetadata(imagePath string, metadata AutoCaptureEmbeddedMetadata) error {
	data, err := os.ReadFile(imagePath) // #nosec G304 -- image path is generated or selected by auto-capture.
	if err != nil {
		return err
	}
	payload, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	switch strings.ToLower(filepath.Ext(imagePath)) {
	case ".png":
		data, err = insertPNGiTXt(data, embeddedMetadataKeyword, payload)
	case ".jpg", ".jpeg":
		data, err = insertJPEGExifDescription(data, payload)
	default:
		return fmt.Errorf("埋め込みメタデータはPNG/JPEGのみ対応です: %s", filepath.Ext(imagePath))
	}
	if err != nil {
		return err
	}
	return WritePrivateFile(imagePath, data)
}

func insertPNGiTXt(data []byte, keyword string, text []byte) ([]byte, error) {
	pngSig := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if !bytes.HasPrefix(data, pngSig) {
		return nil, fmt.Errorf("PNG signatureがありません")
	}
	if len(data) < 33 {
		return nil, fmt.Errorf("PNGデータが短すぎます")
	}
	if string(data[12:16]) != "IHDR" {
		return nil, fmt.Errorf("PNG IHDRが見つかりません")
	}
	insertAt := 8 + 12 + int(binary.BigEndian.Uint32(data[8:12]))
	if insertAt > len(data) {
		return nil, fmt.Errorf("PNG IHDRサイズが不正です")
	}
	var chunkData []byte
	chunkData = append(chunkData, []byte(keyword)...)
	chunkData = append(chunkData, 0)
	chunkData = append(chunkData, 0)
	chunkData = append(chunkData, 0)
	chunkData = append(chunkData, 0)
	chunkData = append(chunkData, 0)
	chunkData = append(chunkData, text...)
	chunk := pngChunk("iTXt", chunkData)
	out := make([]byte, 0, len(data)+len(chunk))
	out = append(out, data[:insertAt]...)
	out = append(out, chunk...)
	out = append(out, data[insertAt:]...)
	return out, nil
}

func pngChunk(kind string, payload []byte) []byte {
	out := make([]byte, 8+len(payload)+4)
	binary.BigEndian.PutUint32(out[0:4], uint32(len(payload)))
	copy(out[4:8], kind)
	copy(out[8:8+len(payload)], payload)
	crc := crc32.ChecksumIEEE(out[4 : 8+len(payload)])
	binary.BigEndian.PutUint32(out[8+len(payload):], crc)
	return out
}

func insertJPEGExifDescription(data []byte, text []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0xff || data[1] != 0xd8 {
		return nil, fmt.Errorf("JPEG SOIがありません")
	}
	exif, err := jpegExifDescriptionSegment(text)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(data)+len(exif))
	out = append(out, data[:2]...)
	out = append(out, exif...)
	out = append(out, data[2:]...)
	return out, nil
}

func jpegExifDescriptionSegment(text []byte) ([]byte, error) {
	if len(text) > 60000 {
		return nil, fmt.Errorf("埋め込みメタデータが大きすぎます")
	}
	value := append(append([]byte{}, text...), 0)
	const ifdOffset = 8
	valueOffset := uint32(ifdOffset + 2 + 12 + 4)
	tiffLen := int(valueOffset) + len(value)
	tiff := make([]byte, tiffLen)
	copy(tiff[0:4], []byte{'I', 'I', 42, 0})
	binary.LittleEndian.PutUint32(tiff[4:8], ifdOffset)
	binary.LittleEndian.PutUint16(tiff[8:10], 1)
	entry := tiff[10:22]
	binary.LittleEndian.PutUint16(entry[0:2], 0x010e)
	binary.LittleEndian.PutUint16(entry[2:4], 2)
	binary.LittleEndian.PutUint32(entry[4:8], uint32(len(value)))
	binary.LittleEndian.PutUint32(entry[8:12], valueOffset)
	copy(tiff[valueOffset:], value)
	payload := append([]byte("Exif\x00\x00"), tiff...)
	if len(payload)+2 > 0xffff {
		return nil, fmt.Errorf("EXIF payloadが大きすぎます")
	}
	segment := make([]byte, 4+len(payload))
	segment[0] = 0xff
	segment[1] = 0xe1
	binary.BigEndian.PutUint16(segment[2:4], uint16(len(payload)+2))
	copy(segment[4:], payload)
	return segment, nil
}
