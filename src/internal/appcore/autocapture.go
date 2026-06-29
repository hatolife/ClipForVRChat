package appcore

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type AutoCaptureEvent struct {
	BatchID string `json:"batchId"`
	ShotID  string `json:"shotId"`
	Path    string `json:"path"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type PresenceUser struct {
	DisplayName string `json:"displayName"`
	UserID      string `json:"userId,omitempty"`
	Status      string `json:"status"`
	Source      string `json:"source"`
	Confidence  string `json:"confidence"`
	JoinedAt    string `json:"joinedAt,omitempty"`
	LeftAt      string `json:"leftAt,omitempty"`
}

type AutoCaptureRunner struct {
	Config  Config
	Handler func(AutoCaptureEvent)
}

func (r AutoCaptureRunner) RunOnce(ctx context.Context) ([]Result, error) {
	cfg := r.Config
	cfg.Normalize()
	ac := cfg.AutoCapture
	if ac.Capture.Mode != "photo" {
		return nil, fmt.Errorf("v0.1.8ではPhoto方式のみ実装済みです")
	}
	views := enabledCameraViews(ac.Views)
	if len(views) == 0 {
		return nil, fmt.Errorf("有効な自動撮影構図がありません")
	}
	batchID := newBatchID(time.Now())
	users, confidence := SnapshotVRChatPresence(ac.Presence.OutputLogDirectory)
	if !ac.Presence.WatchOutputLog {
		users = nil
		confidence = "unknown"
	}
	if !ac.Presence.IncludeUserIDsInSidecar {
		users = presenceUsersWithoutIDs(users)
	}
	photoDir := autoCapturePhotoDirectory(cfg)
	before := scanAutoCapturePhotoFiles(photoDir, ac.Output.Directory)
	client := oscClient{host: ac.OSC.Host, port: ac.OSC.SendPort}
	if err := client.open(); err != nil {
		return nil, err
	}
	defer client.close()
	if err := client.sendInt("/usercamera/Mode", 1); err != nil {
		return nil, err
	}
	if !sleepContext(ctx, 800*time.Millisecond) {
		return nil, ctx.Err()
	}
	results := make([]Result, 0, len(views))
	for i, view := range views {
		if err := ctx.Err(); err != nil {
			return results, err
		}
		shotID := fmt.Sprintf("%s-%02d", batchID, i+1)
		result := r.capturePhotoShot(ctx, client, batchID, shotID, i+1, view, photoDir, before, users, confidence)
		results = append(results, result)
		r.emit(AutoCaptureEvent{BatchID: batchID, ShotID: shotID, Path: result.SourcePath, Error: result.Error, Message: result.Name})
		if result.SourcePath != "" {
			before[result.SourcePath] = time.Now()
		}
	}
	if ac.Capture.CloseCameraAfterBatch {
		_ = client.sendAction("/usercamera/Close")
	}
	return results, nil
}

func (r AutoCaptureRunner) capturePhotoShot(ctx context.Context, client oscClient, batchID string, shotID string, index int, view CameraViewConfig, photoDir string, before map[string]time.Time, users []PresenceUser, confidence string) Result {
	cfg := r.Config.AutoCapture
	name := view.Name
	if name == "" {
		name = view.ID
	}
	if view.CoordinateSpace == "world" || view.Calibrated {
		_ = client.sendFloats("/usercamera/Pose", []float32{
			float32(view.Pose.Position.X), float32(view.Pose.Position.Y), float32(view.Pose.Position.Z),
			float32(view.Pose.Rotation.X), float32(view.Pose.Rotation.Y), float32(view.Pose.Rotation.Z),
		})
	}
	sendOptionalFloat(client, "/usercamera/Zoom", view.Zoom)
	sendOptionalFloat(client, "/usercamera/Exposure", view.Exposure)
	sendOptionalFloat(client, "/usercamera/FocalDistance", view.FocalDistance)
	sendOptionalFloat(client, "/usercamera/Aperture", view.Aperture)
	sendOptionalBool(client, "/usercamera/LookAtMe", view.LookAtMe)
	sendOptionalBool(client, "/usercamera/ShowUIInCamera", view.ShowUIInCamera)
	sendOptionalBool(client, "/usercamera/LocalPlayer", view.LocalPlayer)
	sendOptionalBool(client, "/usercamera/RemotePlayer", view.RemotePlayer)
	sendOptionalBool(client, "/usercamera/Environment", view.Environment)
	settle := time.Duration(cfg.Capture.SettleDelayMS) * time.Millisecond
	if view.SettleDelayMS > 0 {
		settle = time.Duration(view.SettleDelayMS) * time.Millisecond
	}
	if !sleepContext(ctx, settle) {
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	if err := client.sendAction("/usercamera/Capture"); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	if !sleepContext(ctx, time.Duration(cfg.Capture.ButtonReleaseDelayMS)*time.Millisecond) {
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	photoPath := waitForNewPhoto(ctx, photoDir, before, 15*time.Second)
	if photoPath == "" {
		photoPath = waitForNewPhoto(ctx, cfg.Output.Directory, before, 2*time.Second)
	}
	if photoPath == "" {
		return Result{Name: name, Error: "撮影後のVRChat写真ファイルを検出できませんでした。写真保存先設定を確認してください。"}
	}
	if cfg.Output.WriteSidecarJSON {
		_ = WriteAutoCaptureSidecar(photoPath, AutoCaptureSidecar{
			SchemaVersion:   1,
			BatchID:         batchID,
			ShotID:          shotID,
			CapturedAtLocal: time.Now().Format(time.RFC3339),
			CapturedAtUTC:   time.Now().UTC().Format(time.RFC3339),
			CaptureMode:     cfg.Capture.Mode,
			View:            view,
			VRChat: AutoCaptureVRChatMetadata{
				UsersSource:     "output_log",
				UsersConfidence: confidence,
			},
			Users: users,
		})
	}
	result := Result{SourcePath: photoPath, OutputPath: photoPath, Name: filepath.Base(photoPath)}
	if cfg.Discord.Enabled {
		webhook := cfg.Discord.WebhookURL
		if strings.TrimSpace(webhook) == "" {
			webhook = r.Config.Discord.WebhookURL
		}
		uploaded, err := uploadAutoCaptureDiscord(webhook, photoPath, autoCaptureDiscordContent(cfg, view, users))
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.URL = uploaded.URL
		result.DiscordMessageID = uploaded.MessageID
		result.DiscordWebhookID = uploaded.WebhookID
		result.DiscordToken = uploaded.Token
	}
	return result
}

func (r AutoCaptureRunner) emit(event AutoCaptureEvent) {
	if r.Handler != nil {
		r.Handler(event)
	}
}

type AutoCaptureSidecar struct {
	SchemaVersion   int                       `json:"schema_version"`
	BatchID         string                    `json:"batch_id"`
	ShotID          string                    `json:"shot_id"`
	CapturedAtLocal string                    `json:"captured_at_local"`
	CapturedAtUTC   string                    `json:"captured_at_utc"`
	CaptureMode     string                    `json:"capture_mode"`
	View            CameraViewConfig          `json:"view"`
	VRChat          AutoCaptureVRChatMetadata `json:"vrchat"`
	Users           []PresenceUser            `json:"users"`
	Files           AutoCaptureFileMetadata   `json:"files"`
}

type AutoCaptureVRChatMetadata struct {
	WorldID         string `json:"world_id,omitempty"`
	InstanceID      string `json:"instance_id,omitempty"`
	UsersSource     string `json:"users_source"`
	UsersConfidence string `json:"users_confidence"`
}

type AutoCaptureFileMetadata struct {
	ImagePath string `json:"image_path"`
	SHA256    string `json:"sha256"`
}

func WriteAutoCaptureSidecar(imagePath string, sidecar AutoCaptureSidecar) error {
	sum, err := fileSHA256(imagePath)
	if err == nil {
		sidecar.Files.SHA256 = sum
	}
	sidecar.Files.ImagePath = imagePath
	data, err := json.MarshalIndent(sidecar, "", "  ")
	if err != nil {
		return err
	}
	return WritePrivateFile(imagePath+".json", append(data, '\n'))
}

func SnapshotVRChatPresence(logDir string) ([]PresenceUser, string) {
	path := latestVRChatOutputLog(logDir)
	if path == "" {
		return nil, "unknown"
	}
	users, ok := parseVRChatPresenceLog(path)
	if !ok {
		return nil, "unknown"
	}
	out := make([]PresenceUser, 0, len(users))
	for _, user := range users {
		out = append(out, user)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].DisplayName) < strings.ToLower(out[j].DisplayName)
	})
	return out, "partial"
}

func presenceUsersWithoutIDs(users []PresenceUser) []PresenceUser {
	out := make([]PresenceUser, len(users))
	copy(out, users)
	for i := range out {
		out[i].UserID = ""
		if out[i].Confidence == "confirmed" {
			out[i].Confidence = "partial"
		}
	}
	return out
}

func autoCaptureDiscordContent(cfg AutoCaptureConfig, view CameraViewConfig, users []PresenceUser) string {
	lines := []string{
		"VRChat自動撮影",
		"構図: " + view.Name,
		"撮影方式: " + cfg.Capture.Mode,
	}
	if cfg.Presence.IncludeDisplayNamesInDiscord || cfg.Presence.IncludeUserIDsInDiscord {
		parts := make([]string, 0, len(users))
		for _, user := range users {
			name := strings.TrimSpace(user.DisplayName)
			if name == "" {
				name = "unknown"
			}
			if cfg.Presence.IncludeUserIDsInDiscord && user.UserID != "" {
				name += " (" + user.UserID + ")"
			}
			parts = append(parts, name)
		}
		if len(parts) > 0 {
			lines = append(lines, "同席ユーザー: "+strings.Join(parts, ", "))
		}
	}
	return strings.Join(lines, "\n")
}

func uploadAutoCaptureDiscord(webhookURL string, imagePath string, content string) (DiscordUpload, error) {
	var uploaded DiscordUpload
	data, err := os.ReadFile(imagePath) // #nosec G304 -- captured image path comes from configured VRChat photo directory.
	if err != nil {
		return uploaded, err
	}
	ext := strings.ToLower(filepath.Ext(imagePath))
	mime := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".webp":
		mime = "image/webp"
	case ".png":
	default:
		ext = ".png"
	}
	return UploadDiscordWithContent(webhookURL, filepath.Base(imagePath), EncodedImage{Extension: ext, Mime: mime, Data: data}, content)
}

func latestVRChatOutputLog(dir string) string {
	if strings.TrimSpace(dir) == "" {
		dir = DefaultVRChatLogDirectory()
	}
	matches, err := filepath.Glob(filepath.Join(dir, "output_log_*.txt"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	sort.Slice(matches, func(i, j int) bool {
		ii, _ := os.Stat(matches[i])
		jj, _ := os.Stat(matches[j])
		if ii == nil || jj == nil {
			return matches[i] > matches[j]
		}
		return ii.ModTime().After(jj.ModTime())
	})
	return matches[0]
}

var (
	usrIDPattern = regexp.MustCompile(`usr_[0-9a-fA-F-]{36}`)
	joinPattern  = regexp.MustCompile(`(?i)(?:OnPlayerJoined|Joining|joined).*?((?:usr_[0-9a-fA-F-]{36})|$)`)
	leavePattern = regexp.MustCompile(`(?i)(?:OnPlayerLeft|Leaving|left).*?((?:usr_[0-9a-fA-F-]{36})|$)`)
	namePattern  = regexp.MustCompile(`\b(?:displayName|userName|name)[:=]\s*"?([^",\]]+)`)
)

func parseVRChatPresenceLog(path string) (map[string]PresenceUser, bool) {
	data, err := os.ReadFile(path) // #nosec G304 -- VRChat output log path is configured by the local user.
	if err != nil {
		return nil, false
	}
	users := map[string]PresenceUser{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Entering Room") || strings.Contains(line, "Joining wrld_") {
			users = map[string]PresenceUser{}
			continue
		}
		id := usrIDPattern.FindString(line)
		name := extractPresenceName(line)
		switch {
		case joinPattern.MatchString(line):
			key := id
			if key == "" {
				key = name
			}
			if key != "" {
				users[key] = PresenceUser{DisplayName: name, UserID: id, Status: "present", Source: "output_log", Confidence: presenceConfidence(id), JoinedAt: extractLogTime(line)}
			}
		case leavePattern.MatchString(line):
			key := id
			if key == "" {
				key = name
			}
			delete(users, key)
		}
	}
	return users, true
}

func extractPresenceName(line string) string {
	if match := namePattern.FindStringSubmatch(line); len(match) == 2 {
		name := strings.TrimSpace(match[1])
		if idx := strings.Index(name, " usr_"); idx >= 0 {
			name = strings.TrimSpace(name[:idx])
		}
		return name
	}
	if idx := strings.LastIndex(line, "usr_"); idx > 0 {
		prefix := strings.TrimSpace(line[:idx])
		fields := strings.Fields(prefix)
		if len(fields) > 0 {
			return strings.Trim(fields[len(fields)-1], `"'[]():`)
		}
	}
	return ""
}

func presenceConfidence(id string) string {
	if id == "" {
		return "partial"
	}
	return "confirmed"
}

func extractLogTime(line string) string {
	if len(line) >= 19 {
		prefix := line[:19]
		if _, err := time.Parse("2006.01.02 15:04:05", prefix); err == nil {
			return prefix
		}
	}
	return ""
}

func enabledCameraViews(views []CameraViewConfig) []CameraViewConfig {
	out := make([]CameraViewConfig, 0, len(views))
	for _, view := range views {
		if view.Enabled {
			out = append(out, view)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].SortOrder < out[j].SortOrder
	})
	return out
}

func waitForNewPhoto(ctx context.Context, dir string, before map[string]time.Time, timeout time.Duration) string {
	if strings.TrimSpace(dir) == "" {
		return ""
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		current := scanPhotoFiles(dir)
		paths := sortedPhotoPaths(current)
		for i := len(paths) - 1; i >= 0; i-- {
			path := paths[i]
			if _, ok := before[path]; ok {
				continue
			}
			if fileLooksStable(path) {
				return path
			}
		}
		if !sleepContext(ctx, 500*time.Millisecond) {
			return ""
		}
	}
	return ""
}

func scanAutoCapturePhotoFiles(photoDir string, outputDir string) map[string]time.Time {
	files := scanPhotoFiles(photoDir)
	for path, modTime := range scanPhotoFiles(outputDir) {
		files[path] = modTime
	}
	return files
}

func autoCapturePhotoDirectory(cfg Config) string {
	photoDir := strings.TrimSpace(cfg.AutoPhoto.PhotoDirectory)
	if photoDir == "" {
		return DefaultVRChatPhotoDirectory()
	}
	return photoDir
}

func newBatchID(t time.Time) string {
	return "batch-" + t.Format("20060102-150405")
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- hashing the captured local image selected by the app.
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func sendOptionalFloat(client oscClient, address string, value *float64) {
	if value != nil {
		_ = client.sendFloat(address, float32(*value))
	}
}

func sendOptionalBool(client oscClient, address string, value *bool) {
	if value != nil {
		_ = client.sendBool(address, *value)
	}
}

type oscClient struct {
	host string
	port int
	conn net.Conn
}

func (c *oscClient) open() error {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", c.host, c.port), 3*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *oscClient) close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c oscClient) sendInt(address string, value int32) error {
	return c.send(address, ",i", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], uint32(value))
		return append(buf, raw[:]...)
	})
}

func (c oscClient) sendFloat(address string, value float32) error {
	return c.send(address, ",f", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], math.Float32bits(value))
		return append(buf, raw[:]...)
	})
}

func (c oscClient) sendFloats(address string, values []float32) error {
	types := "," + strings.Repeat("f", len(values))
	return c.send(address, types, func(buf []byte) []byte {
		for _, value := range values {
			var raw [4]byte
			binary.BigEndian.PutUint32(raw[:], math.Float32bits(value))
			buf = append(buf, raw[:]...)
		}
		return buf
	})
}

func (c oscClient) sendBool(address string, value bool) error {
	tag := ",F"
	if value {
		tag = ",T"
	}
	return c.send(address, tag, func(buf []byte) []byte { return buf })
}

func (c oscClient) sendAction(address string) error {
	return c.send(address, ",", func(buf []byte) []byte { return buf })
}

func (c oscClient) send(address string, typeTags string, appendArgs func([]byte) []byte) error {
	if c.conn == nil {
		return fmt.Errorf("OSC接続が開かれていません")
	}
	_, err := c.conn.Write(buildOSCPacket(address, typeTags, appendArgs))
	return err
}

func buildOSCPacket(address string, typeTags string, appendArgs func([]byte) []byte) []byte {
	packet := appendOSCString(nil, address)
	packet = appendOSCString(packet, typeTags)
	return appendArgs(packet)
}

func appendOSCString(buf []byte, value string) []byte {
	buf = append(buf, []byte(value)...)
	buf = append(buf, 0)
	for len(buf)%4 != 0 {
		buf = append(buf, 0)
	}
	return buf
}
