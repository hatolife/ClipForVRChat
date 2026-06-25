package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

const releaseSigningPublicKeyArmored = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Comment: A256 3716 E502 15E2 5B75  E4EE F603 3F2A 3D70 179C
Comment: hatolife <poppo@hato.life>

xjMEaje/QhYJKwYBBAHaRw8BAQdAYlhl12X7AOD3rzzIf+0FjzIXNg9SVTJg12FU
qJMcH43NGmhhdG9saWZlIDxwb3Bwb0BoYXRvLmxpZmU+wpMEExYKADsWIQSiVjcW
5QIV4lt15O72Az8qPXAXnAUCaje/QgIbAwULCQgHAgIiAgYVCgkICwIEFgIDAQIe
BwIXgAAKCRD2Az8qPXAXnCEOAQCcXCANFkwjq/adDUSGYMBQkiqypymQccozJYkN
L07bigD8CfBRI3wCInG5ThrCE7rxxxwfqvRQklXGpDSdwTAaUQfOOARqN79CEgor
BgEEAZdVAQUBAQdAw/TLMagbQzTDepQPqC3oSJK7JFGPLS6onLPPnVXCvwoDAQgH
wngEGBYKACAWIQSiVjcW5QIV4lt15O72Az8qPXAXnAUCaje/QgIbDAAKCRD2Az8q
PXAXnFGqAQDhmF0+rCS2r4Ya1uNS1zAujd0JL3pmBscSUK4EavGd5wEAibO6Vwd1
mwM/zreTlthl15c9e1qSIEcR5OdubzCuJgI=
=OlRJ
-----END PGP PUBLIC KEY BLOCK-----`

type diagnosticPackageManifest struct {
	CreatedAt    string                  `json:"createdAt"`
	AppVersion   string                  `json:"appVersion"`
	Version      string                  `json:"version"`
	Revision     string                  `json:"revision"`
	ReleaseTime  string                  `json:"releaseTime"`
	Executable   string                  `json:"executable"`
	ExeSHA256    string                  `json:"exeSha256"`
	ConfigPath   string                  `json:"configPath"`
	HistoryPath  string                  `json:"historyPath"`
	LogPath      string                  `json:"logPath"`
	Config       json.RawMessage         `json:"configSummary"`
	Files        []diagnosticPackageFile `json:"files"`
	MissingFiles []string                `json:"missingFiles"`
}

type diagnosticPackageFile struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256,omitempty"`
	Size   int64  `json:"size"`
}

type diagnosticPathRedactor struct {
	replacements []string
}

type diagnosticPathReplacement struct {
	source      string
	placeholder string
}

func createEncryptedDiagnosticPackage(configPath string, cfg appcore.Config) (string, error) {
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(releaseSigningPublicKeyArmored))
	if err != nil {
		return "", fmt.Errorf("診断パッケージ用公開鍵を読み込めません: %w", err)
	}
	if !hasEncryptionCapableKey(entities) {
		return "", fmt.Errorf("診断パッケージ用公開鍵は暗号化用途に対応していません。暗号化用の公開鍵を同梱する必要があります")
	}

	outputDir := filepath.Dir(configPath)
	if strings.TrimSpace(configPath) == "" {
		outputDir = "."
	}
	stamp := time.Now().Format("20060102-150405")
	workDir := filepath.Join(outputDir, "diagnostics", stamp)
	dataDir := filepath.Join(workDir, "data")
	zipPath := filepath.Join(workDir, fmt.Sprintf("ClipForVRChat-diagnostics-%s.zip", stamp))
	outputPath := zipPath + ".gpg"
	redactor := newDiagnosticPathRedactor(configPath)

	appendDiagnosticOutputDirectoryLog(configPath, cfg, redactor)

	manifest, err := prepareDiagnosticDataDirectory(configPath, cfg, dataDir)
	if err != nil {
		return "", err
	}
	if err := zipDirectory(dataDir, zipPath); err != nil {
		return "", err
	}
	if err := os.RemoveAll(dataDir); err != nil {
		return "", fmt.Errorf("診断データ一時フォルダを削除できません: %w", err)
	}
	zipData, err := os.ReadFile(zipPath) // #nosec G304 -- zip path is generated inside the app-owned diagnostics work directory.
	if err != nil {
		return "", fmt.Errorf("診断zipを読み込めません: %w", err)
	}
	encrypted, err := encryptDiagnosticZip(zipData, entities)
	if err != nil {
		return "", err
	}
	if err := appcore.WritePrivateFile(outputPath, encrypted); err != nil {
		return "", fmt.Errorf("診断パッケージを保存できません: %w", err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "diagnostic package=%q zip=%q staged_data_dir=%q files=%d missing=%d", outputPath, zipPath, dataDir, len(manifest.Files), len(manifest.MissingFiles))
	return outputPath, nil
}

func encryptZipFileWithPublicKey(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("zipファイルを指定してください")
	}
	if !strings.EqualFold(filepath.Ext(path), ".zip") {
		return "", fmt.Errorf("zipファイルではありません: %s", path)
	}
	info, err := os.Stat(path) // #nosec G304,G703 -- path comes from explicit CLI zip argument and is validated as a local .zip file.
	if err != nil {
		return "", fmt.Errorf("zipファイルを確認できません: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("zipファイルではありません: %s", path)
	}
	data, err := os.ReadFile(path) // #nosec G304,G703 -- path comes from explicit CLI zip argument and is validated as a local .zip file.
	if err != nil {
		return "", fmt.Errorf("zipファイルを読み込めません: %w", err)
	}
	if _, err := zip.NewReader(bytes.NewReader(data), int64(len(data))); err != nil {
		return "", fmt.Errorf("zipファイルを読み込めません: %w", err)
	}
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(releaseSigningPublicKeyArmored))
	if err != nil {
		return "", fmt.Errorf("公開鍵を読み込めません: %w", err)
	}
	if !hasEncryptionCapableKey(entities) {
		return "", fmt.Errorf("公開鍵は暗号化用途に対応していません")
	}
	encrypted, err := encryptDiagnosticZip(data, entities)
	if err != nil {
		return "", err
	}
	outputPath := path + ".gpg"
	if err := appcore.WritePrivateFile(outputPath, encrypted); err != nil {
		return "", fmt.Errorf("暗号化zipを保存できません: %w", err)
	}
	return outputPath, nil
}

func encryptDiagnosticZip(zipData []byte, entities openpgp.EntityList) ([]byte, error) {
	var encrypted bytes.Buffer
	writer, err := openpgp.Encrypt(&encrypted, entities, nil, &openpgp.FileHints{
		IsBinary: true,
		FileName: "diagnostics.zip",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("診断パッケージを暗号化できません: %w", err)
	}
	if _, err := writer.Write(zipData); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("診断パッケージを書き込めません: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("診断パッケージを完了できません: %w", err)
	}
	return encrypted.Bytes(), nil
}

func hasEncryptionCapableKey(entities openpgp.EntityList) bool {
	for _, entity := range entities {
		if len(entities.KeysByIdUsage(entity.PrimaryKey.KeyId, packet.KeyFlagEncryptCommunications)) > 0 ||
			len(entities.KeysByIdUsage(entity.PrimaryKey.KeyId, packet.KeyFlagEncryptStorage)) > 0 {
			return true
		}
		for _, subkey := range entity.Subkeys {
			if len(entities.KeysByIdUsage(subkey.PublicKey.KeyId, packet.KeyFlagEncryptCommunications)) > 0 ||
				len(entities.KeysByIdUsage(subkey.PublicKey.KeyId, packet.KeyFlagEncryptStorage)) > 0 {
				return true
			}
		}
	}
	return false
}

func buildDiagnosticZip(configPath string, cfg appcore.Config) ([]byte, diagnosticPackageManifest, error) {
	tempDir, err := os.MkdirTemp("", "clipforvrchat-diagnostics-*")
	if err != nil {
		return nil, diagnosticPackageManifest{}, err
	}
	defer os.RemoveAll(tempDir)
	dataDir := filepath.Join(tempDir, "data")
	zipPath := filepath.Join(tempDir, "diagnostics.zip")
	manifest, err := prepareDiagnosticDataDirectory(configPath, cfg, dataDir)
	if err != nil {
		return nil, manifest, err
	}
	if err := zipDirectory(dataDir, zipPath); err != nil {
		return nil, manifest, err
	}
	data, err := os.ReadFile(zipPath) // #nosec G304 -- zip path is generated inside a temporary diagnostics work directory.
	if err != nil {
		return nil, manifest, err
	}
	return data, manifest, nil
}

func prepareDiagnosticDataDirectory(configPath string, cfg appcore.Config, dataDir string) (diagnosticPackageManifest, error) {
	exePath, err := os.Executable()
	if err != nil {
		exePath = ""
	}
	exeSHA256 := ""
	if exePath != "" {
		exeSHA256, _ = fileSHA256(exePath)
	}
	configJSON := json.RawMessage(configSummaryForLog(cfg))
	redactor := newDiagnosticPathRedactor(configPath)
	manifest := diagnosticPackageManifest{
		CreatedAt:   time.Now().Format(time.RFC3339),
		AppVersion:  appVersion(),
		Version:     version,
		Revision:    revision,
		ReleaseTime: appReleaseTime(),
		Executable:  redactor.Redact(exePath),
		ExeSHA256:   exeSHA256,
		ConfigPath:  redactor.Redact(configPath),
		HistoryPath: redactor.Redact(appcore.HistoryPath(configPath)),
		LogPath:     redactor.Redact(appcore.DiagnosticLogPath(configPath)),
		Config:      json.RawMessage(redactor.Redact(string(configJSON))),
	}

	add := func(name, path string) error {
		info, err := os.Stat(path)
		if err != nil {
			manifest.MissingFiles = append(manifest.MissingFiles, name)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		hash, _ := fileSHA256(path)
		manifest.Files = append(manifest.Files, diagnosticPackageFile{Name: name, Path: redactor.Redact(path), SHA256: hash, Size: info.Size()})
		if isDiagnosticTextFile(name) {
			return addRedactedTextFileToDirectory(dataDir, name, path, redactor)
		}
		return copyFileToDirectory(dataDir, name, path)
	}

	if err := add("config.json", configPath); err != nil {
		return manifest, err
	}
	if err := add("history.json", appcore.HistoryPath(configPath)); err != nil {
		return manifest, err
	}
	if exePath != "" {
		if err := add("ClipForVRChat.exe", exePath); err != nil {
			return manifest, err
		}
	} else {
		manifest.MissingFiles = append(manifest.MissingFiles, "ClipForVRChat.exe")
	}
	logFiles, _ := filepath.Glob(filepath.Join(appcore.DiagnosticLogDir(configPath), "*.log"))
	for _, logPath := range logFiles {
		name := filepath.Join("logs", filepath.Base(logPath))
		if err := add(name, logPath); err != nil {
			return manifest, err
		}
	}
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return manifest, err
	}
	if err := writePrivateFileInDirectory(dataDir, "manifest.json", append(manifestData, '\n')); err != nil {
		return manifest, err
	}
	return manifest, nil
}

func appendDiagnosticOutputDirectoryLog(configPath string, cfg appcore.Config, redactor diagnosticPathRedactor) {
	logPath := appcore.DiagnosticLogPath(configPath)
	outputDir := diagnosticOutputDirectory(configPath, cfg)
	if outputDir == "" {
		appcore.AppendDiagnosticLog(logPath, "diagnostic output directory is not configured")
		return
	}
	info, err := os.Stat(outputDir) // #nosec G304 -- output directory path comes from user configuration and is only listed for diagnostics.
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "diagnostic output directory=%q status=missing error=%q", redactor.Redact(outputDir), err.Error())
		return
	}
	if !info.IsDir() {
		appcore.AppendDiagnosticLog(logPath, "diagnostic output directory=%q status=not_directory", redactor.Redact(outputDir))
		return
	}
	entries, err := os.ReadDir(outputDir) // #nosec G304 -- output directory path comes from user configuration and is only listed for diagnostics.
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "diagnostic output directory=%q status=read_error error=%q", redactor.Redact(outputDir), err.Error())
		return
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	appcore.AppendDiagnosticLog(logPath, "diagnostic output directory=%q entries=%d", redactor.Redact(outputDir), len(entries))
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			appcore.AppendDiagnosticLog(logPath, "diagnostic output entry name=%q status=stat_error error=%q", entry.Name(), err.Error())
			continue
		}
		appcore.AppendDiagnosticLog(logPath, "diagnostic output entry name=%q size=%d dir=%t", entry.Name(), entryInfo.Size(), entryInfo.IsDir())
	}
}

func diagnosticOutputDirectory(configPath string, cfg appcore.Config) string {
	baseDir := filepath.Dir(configPath)
	outputDir := strings.TrimSpace(cfg.Image.OutputDirectory)
	if outputDir == "" {
		outputDir = "./output"
	}
	if !filepath.IsAbs(outputDir) && strings.TrimSpace(baseDir) != "" {
		outputDir = filepath.Join(baseDir, outputDir)
	}
	return outputDir
}

func copyFileToDirectory(baseDir string, name string, path string) error {
	file, err := os.Open(path) // #nosec G304 -- diagnostic package paths are derived from app-owned config/history/log paths and the running executable.
	if err != nil {
		return fmt.Errorf("%s を読み込めません: %w", name, err)
	}
	defer file.Close()
	target := filepath.Join(baseDir, filepath.FromSlash(filepath.ToSlash(name)))
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return err
	}
	writer, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600) // #nosec G304 -- target path is generated under the diagnostics work directory from fixed package entry names.
	if err != nil {
		return fmt.Errorf("%s を作成できません: %w", name, err)
	}
	defer writer.Close()
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("%s を書き込めません: %w", name, err)
	}
	return nil
}

func addRedactedTextFileToDirectory(baseDir string, name string, path string, redactor diagnosticPathRedactor) error {
	data, err := os.ReadFile(path) // #nosec G304 -- diagnostic text paths are app-owned config/history/log files.
	if err != nil {
		return fmt.Errorf("%s を読み込めません: %w", name, err)
	}
	return writePrivateFileInDirectory(baseDir, name, []byte(redactor.Redact(string(data))))
}

func writePrivateFileInDirectory(baseDir string, name string, data []byte) error {
	target := filepath.Join(baseDir, filepath.FromSlash(filepath.ToSlash(name)))
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return err
	}
	return appcore.WritePrivateFile(target, data)
}

func zipDirectory(sourceDir string, zipPath string) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	err := filepath.WalkDir(sourceDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if err := addFileToZip(zipWriter, filepath.ToSlash(rel), path); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		_ = zipWriter.Close()
		return err
	}
	if err := zipWriter.Close(); err != nil {
		return err
	}
	return appcore.WritePrivateFile(zipPath, buf.Bytes())
}

func addFileToZip(zipWriter *zip.Writer, name, path string) error {
	file, err := os.Open(path) // #nosec G304 -- zip source paths are generated in diagnostics work directory.
	if err != nil {
		return fmt.Errorf("%s を読み込めません: %w", name, err)
	}
	defer file.Close()
	header := &zip.FileHeader{Name: filepath.ToSlash(name), Method: zip.Deflate}
	header.SetMode(0600)
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("%s をzipに追加できません: %w", name, err)
	}
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("%s をzipへ書き込めません: %w", name, err)
	}
	return nil
}

func isDiagnosticTextFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".json", ".log", ".txt":
		return true
	default:
		return false
	}
}

func newDiagnosticPathRedactor(configPath string) diagnosticPathRedactor {
	candidates := []struct {
		name  string
		value string
	}{
		{"USERPROFILE", os.Getenv("USERPROFILE")},
		{"APPDATA", os.Getenv("APPDATA")},
		{"LOCALAPPDATA", os.Getenv("LOCALAPPDATA")},
		{"TEMP", os.Getenv("TEMP")},
		{"TMP", os.Getenv("TMP")},
		{"HOME", os.Getenv("HOME")},
	}
	if strings.TrimSpace(configPath) != "" {
		candidates = append(candidates, struct {
			name  string
			value string
		}{"CLIPFORVRCHAT_DIR", filepath.Dir(configPath)})
	}
	pairs := make([]diagnosticPathReplacement, 0, len(candidates)*4)
	seen := map[string]bool{}
	for _, candidate := range candidates {
		value := strings.Trim(strings.TrimSpace(candidate.value), `"`)
		if value == "" {
			continue
		}
		addPathReplacement(&pairs, seen, value, "%"+candidate.name+"%")
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		return len(pairs[i].source) > len(pairs[j].source)
	})
	replacements := make([]string, 0, len(pairs)*2)
	for _, pair := range pairs {
		replacements = append(replacements, pair.source, pair.placeholder)
	}
	return diagnosticPathRedactor{replacements: replacements}
}

func addPathReplacement(replacements *[]diagnosticPathReplacement, seen map[string]bool, path string, placeholder string) {
	variants := []string{filepath.Clean(path)}
	slash := strings.ReplaceAll(filepath.Clean(path), `\`, `/`)
	backslash := strings.ReplaceAll(filepath.Clean(path), `/`, `\`)
	escapedBackslash := strings.ReplaceAll(backslash, `\`, `\\`)
	variants = append(variants, slash, backslash, escapedBackslash)
	sort.SliceStable(variants, func(i, j int) bool {
		return len(variants[i]) > len(variants[j])
	})
	for _, value := range variants {
		if value == "." || value == "" || seen[value] {
			continue
		}
		seen[value] = true
		*replacements = append(*replacements, diagnosticPathReplacement{source: value, placeholder: placeholder})
	}
}

func (r diagnosticPathRedactor) Redact(text string) string {
	if text == "" || len(r.replacements) == 0 {
		return text
	}
	replacer := strings.NewReplacer(r.replacements...)
	return replacer.Replace(text)
}
