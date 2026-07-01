package appcore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"time"
)

const (
	embeddedMetadataKeyword       = "ClipForVRChat:AutoCapture"
	embeddedMetadataApp           = "ClipForVRChat"
	embeddedMetadataAppVersion    = "v0.1.8"
	maxEmbeddedMetadataPayloadLen = 60000
)

var pngSignature = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}

type AutoCaptureEmbeddedMetadata struct {
	SchemaVersion      int                        `json:"schema_version"`
	App                string                     `json:"app"`
	AppVersion         string                     `json:"app_version"`
	BatchID            string                     `json:"batch_id"`
	ShotID             string                     `json:"shot_id"`
	CapturedAtUTC      string                     `json:"captured_at_utc"`
	CaptureMode        string                     `json:"capture_mode"`
	View               CameraViewConfig           `json:"view"`
	ResolvedPose       *CameraPoseConfig          `json:"resolved_pose,omitempty"`
	Stream             *AutoCaptureStreamMetadata `json:"stream,omitempty"`
	PresenceSource     string                     `json:"presence_source"`
	PresenceConfidence string                     `json:"presence_confidence"`
	UserCount          int                        `json:"user_count"`
	UsersTruncated     bool                       `json:"users_truncated,omitempty"`
	MetadataWarnings   []string                   `json:"metadata_warnings,omitempty"`
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
		App:                embeddedMetadataApp,
		AppVersion:         embeddedMetadataAppVersion,
		BatchID:            batchID,
		ShotID:             shotID,
		CapturedAtUTC:      time.Now().UTC().Format(time.RFC3339),
		CaptureMode:        cfg.Capture.Mode,
		View:               view,
		Stream:             autoCaptureStreamMetadata(streamInfo),
		PresenceSource:     "output_log",
		PresenceConfidence: confidence,
		UserCount:          len(metadataUsers),
		Users:              metadataUsers,
	}
}

func WriteAutoCaptureEmbeddedMetadata(imagePath string, metadata AutoCaptureEmbeddedMetadata) error {
	_, err := WriteAutoCaptureEmbeddedMetadataWithWarnings(imagePath, metadata)
	return err
}

func WriteAutoCaptureEmbeddedMetadataWithWarnings(imagePath string, metadata AutoCaptureEmbeddedMetadata) ([]string, error) {
	data, err := os.ReadFile(imagePath) // #nosec G304 -- image path is generated or selected by auto-capture.
	if err != nil {
		return nil, err
	}
	payload, metadata, err := marshalEmbeddedMetadataPayload(metadata)
	if err != nil {
		return nil, err
	}
	switch detectEmbeddedImageFormat(data) {
	case "png":
		data, err = insertPNGEmbeddedMetadata(data, payload)
	case "jpeg":
		data, err = insertJPEGEmbeddedMetadata(data, payload)
	default:
		return metadata.MetadataWarnings, fmt.Errorf("埋め込みメタデータはPNG/JPEGのみ対応です")
	}
	if err != nil {
		return metadata.MetadataWarnings, err
	}
	return metadata.MetadataWarnings, WritePrivateFile(imagePath, data)
}

func ReadAutoCaptureEmbeddedMetadata(imagePath string) (AutoCaptureEmbeddedMetadata, error) {
	data, err := os.ReadFile(imagePath) // #nosec G304 -- caller supplied local image path.
	if err != nil {
		return AutoCaptureEmbeddedMetadata{}, err
	}
	var payload []byte
	switch detectEmbeddedImageFormat(data) {
	case "png":
		payload, err = readPNGEmbeddedMetadata(data)
	case "jpeg":
		payload, err = readJPEGEmbeddedMetadata(data)
	default:
		err = fmt.Errorf("埋め込みメタデータはPNG/JPEGのみ対応です")
	}
	if err != nil {
		return AutoCaptureEmbeddedMetadata{}, err
	}
	var metadata AutoCaptureEmbeddedMetadata
	if err := json.Unmarshal(payload, &metadata); err != nil {
		return AutoCaptureEmbeddedMetadata{}, err
	}
	return metadata, nil
}

func marshalEmbeddedMetadataPayload(metadata AutoCaptureEmbeddedMetadata) ([]byte, AutoCaptureEmbeddedMetadata, error) {
	metadata.App = embeddedMetadataApp
	if metadata.AppVersion == "" {
		metadata.AppVersion = embeddedMetadataAppVersion
	}
	metadata.UserCount = len(metadata.Users)
	payload, err := json.Marshal(metadata)
	if err != nil {
		return nil, metadata, err
	}
	if len(payload) <= maxEmbeddedMetadataPayloadLen {
		return payload, metadata, nil
	}
	originalCount := len(metadata.Users)
	truncated := false
	for len(metadata.Users) > 0 && len(payload) > maxEmbeddedMetadataPayloadLen {
		metadata.Users = metadata.Users[:len(metadata.Users)-1]
		metadata.UsersTruncated = true
		truncated = true
		metadata.UserCount = len(metadata.Users)
		payload, err = json.Marshal(metadata)
		if err != nil {
			return nil, metadata, err
		}
	}
	if truncated {
		metadata.MetadataWarnings = appendUniqueString(metadata.MetadataWarnings, fmt.Sprintf("metadata users were truncated from %d to %d to fit embedded payload limit", originalCount, len(metadata.Users)))
		payload, err = json.Marshal(metadata)
		if err != nil {
			return nil, metadata, err
		}
		for len(metadata.Users) > 0 && len(payload) > maxEmbeddedMetadataPayloadLen {
			metadata.Users = metadata.Users[:len(metadata.Users)-1]
			metadata.UserCount = len(metadata.Users)
			payload, err = json.Marshal(metadata)
			if err != nil {
				return nil, metadata, err
			}
		}
	}
	if len(payload) > maxEmbeddedMetadataPayloadLen {
		return nil, metadata, fmt.Errorf("埋め込みメタデータが大きすぎます")
	}
	return payload, metadata, nil
}

func detectEmbeddedImageFormat(data []byte) string {
	switch {
	case bytes.HasPrefix(data, pngSignature):
		return "png"
	case len(data) >= 2 && data[0] == 0xff && data[1] == 0xd8:
		return "jpeg"
	default:
		return ""
	}
}

func insertPNGEmbeddedMetadata(data []byte, payload []byte) ([]byte, error) {
	cleaned, insertAt, err := removePNGOwnMetadata(data)
	if err != nil {
		return nil, err
	}
	itxt := pngChunk("iTXt", pngITXtPayload(embeddedMetadataKeyword, payload))
	exif := pngChunk("eXIf", tiffImageDescription(payload))
	out := make([]byte, 0, len(cleaned)+len(itxt)+len(exif))
	out = append(out, cleaned[:insertAt]...)
	out = append(out, itxt...)
	out = append(out, exif...)
	out = append(out, cleaned[insertAt:]...)
	return out, nil
}

func readPNGEmbeddedMetadata(data []byte) ([]byte, error) {
	if !bytes.HasPrefix(data, pngSignature) {
		return nil, fmt.Errorf("PNG signatureがありません")
	}
	var exifPayload []byte
	for offset := 8; offset+12 <= len(data); {
		size := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		if size < 0 || offset+12+size > len(data) {
			return nil, fmt.Errorf("PNG chunkサイズが不正です")
		}
		kind := string(data[offset+4 : offset+8])
		payload := data[offset+8 : offset+8+size]
		if kind == "iTXt" {
			if text, ok := parsePNGiTXtPayload(payload, embeddedMetadataKeyword); ok {
				return text, nil
			}
		}
		if kind == "eXIf" && bytes.Contains(payload, []byte(embeddedMetadataApp)) {
			exifPayload = payload
		}
		offset += 12 + size
	}
	if exifPayload != nil {
		return readTIFFImageDescription(exifPayload)
	}
	return nil, fmt.Errorf("自動撮影メタデータが見つかりません")
}

func removePNGOwnMetadata(data []byte) ([]byte, int, error) {
	if !bytes.HasPrefix(data, pngSignature) {
		return nil, 0, fmt.Errorf("PNG signatureがありません")
	}
	out := append([]byte{}, data[:8]...)
	insertAt := 0
	for offset := 8; offset+12 <= len(data); {
		size := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		if size < 0 || offset+12+size > len(data) {
			return nil, 0, fmt.Errorf("PNG chunkサイズが不正です")
		}
		kind := string(data[offset+4 : offset+8])
		payload := data[offset+8 : offset+8+size]
		chunk := data[offset : offset+12+size]
		remove := false
		if kind == "iTXt" {
			_, remove = parsePNGiTXtPayload(payload, embeddedMetadataKeyword)
		}
		if kind == "eXIf" && bytes.Contains(payload, []byte(embeddedMetadataApp)) {
			remove = true
		}
		if !remove {
			out = append(out, chunk...)
			if kind == "IHDR" {
				insertAt = len(out)
			}
		}
		offset += 12 + size
	}
	if insertAt == 0 {
		return nil, 0, fmt.Errorf("PNG IHDRが見つかりません")
	}
	return out, insertAt, nil
}

func pngITXtPayload(keyword string, text []byte) []byte {
	var chunkData []byte
	chunkData = append(chunkData, []byte(keyword)...)
	chunkData = append(chunkData, 0, 0, 0, 0, 0)
	chunkData = append(chunkData, text...)
	return chunkData
}

func parsePNGiTXtPayload(payload []byte, keyword string) ([]byte, bool) {
	prefix := append([]byte(keyword), 0, 0, 0, 0, 0)
	if !bytes.HasPrefix(payload, prefix) {
		return nil, false
	}
	return payload[len(prefix):], true
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

func insertJPEGEmbeddedMetadata(data []byte, payload []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0xff || data[1] != 0xd8 {
		return nil, fmt.Errorf("JPEG SOIがありません")
	}
	cleaned, insertAt, err := removeJPEGOwnMetadata(data)
	if err != nil {
		return nil, err
	}
	exif, err := jpegExifDescriptionSegment(payload)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(cleaned)+len(exif))
	out = append(out, cleaned[:insertAt]...)
	out = append(out, exif...)
	out = append(out, cleaned[insertAt:]...)
	return out, nil
}

func readJPEGEmbeddedMetadata(data []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0xff || data[1] != 0xd8 {
		return nil, fmt.Errorf("JPEG SOIがありません")
	}
	for offset := 2; offset+4 <= len(data); {
		if data[offset] != 0xff {
			break
		}
		marker := data[offset+1]
		if marker == 0xda || marker == 0xd9 {
			break
		}
		size := int(binary.BigEndian.Uint16(data[offset+2 : offset+4]))
		if size < 2 || offset+2+size > len(data) {
			return nil, fmt.Errorf("JPEG segmentサイズが不正です")
		}
		segment := data[offset+4 : offset+2+size]
		if marker == 0xe1 && bytes.HasPrefix(segment, []byte("Exif\x00\x00")) && bytes.Contains(segment, []byte(embeddedMetadataApp)) {
			return readTIFFImageDescription(segment[6:])
		}
		offset += 2 + size
	}
	return nil, fmt.Errorf("自動撮影メタデータが見つかりません")
}

func removeJPEGOwnMetadata(data []byte) ([]byte, int, error) {
	out := append([]byte{}, data[:2]...)
	insertAt := 2
	offset := 2
	for offset+4 <= len(data) {
		if data[offset] != 0xff {
			break
		}
		marker := data[offset+1]
		if marker == 0xda || marker == 0xd9 {
			break
		}
		size := int(binary.BigEndian.Uint16(data[offset+2 : offset+4]))
		if size < 2 || offset+2+size > len(data) {
			return nil, 0, fmt.Errorf("JPEG segmentサイズが不正です")
		}
		segment := data[offset : offset+2+size]
		payload := data[offset+4 : offset+2+size]
		remove := marker == 0xe1 && bytes.HasPrefix(payload, []byte("Exif\x00\x00")) && bytes.Contains(payload, []byte(embeddedMetadataApp))
		if !remove {
			out = append(out, segment...)
			insertAt = len(out)
		}
		offset += 2 + size
	}
	out = append(out, data[offset:]...)
	return out, insertAt, nil
}

func jpegExifDescriptionSegment(text []byte) ([]byte, error) {
	if len(text) > maxEmbeddedMetadataPayloadLen {
		return nil, fmt.Errorf("EXIF payloadが大きすぎます")
	}
	tiff := tiffImageDescription(text)
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

func tiffImageDescription(text []byte) []byte {
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
	return tiff
}

func readTIFFImageDescription(tiff []byte) ([]byte, error) {
	if len(tiff) < 14 || string(tiff[:2]) != "II" || binary.LittleEndian.Uint16(tiff[2:4]) != 42 {
		return nil, fmt.Errorf("TIFFヘッダーが不正です")
	}
	ifdOffset := int(binary.LittleEndian.Uint32(tiff[4:8]))
	if ifdOffset < 0 || ifdOffset+2 > len(tiff) {
		return nil, fmt.Errorf("TIFF IFD offsetが不正です")
	}
	count := int(binary.LittleEndian.Uint16(tiff[ifdOffset : ifdOffset+2]))
	for i := 0; i < count; i++ {
		entryOffset := ifdOffset + 2 + i*12
		if entryOffset+12 > len(tiff) {
			return nil, fmt.Errorf("TIFF IFD entryが不正です")
		}
		tag := binary.LittleEndian.Uint16(tiff[entryOffset : entryOffset+2])
		fieldType := binary.LittleEndian.Uint16(tiff[entryOffset+2 : entryOffset+4])
		valueCount := int(binary.LittleEndian.Uint32(tiff[entryOffset+4 : entryOffset+8]))
		valueOffset := int(binary.LittleEndian.Uint32(tiff[entryOffset+8 : entryOffset+12]))
		if tag != 0x010e || fieldType != 2 {
			continue
		}
		if valueCount <= 0 || valueOffset < 0 || valueOffset+valueCount > len(tiff) {
			return nil, fmt.Errorf("TIFF ImageDescriptionが不正です")
		}
		value := tiff[valueOffset : valueOffset+valueCount]
		return bytes.TrimRight(value, "\x00"), nil
	}
	return nil, fmt.Errorf("TIFF ImageDescriptionが見つかりません")
}

func appendUniqueString(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
