package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

func createEncryptedDiagnosticPackage(configPath string, cfg appcore.Config) (string, error) {
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(releaseSigningPublicKeyArmored))
	if err != nil {
		return "", fmt.Errorf("診断パッケージ用公開鍵を読み込めません: %w", err)
	}
	if !hasEncryptionCapableKey(entities) {
		return "", fmt.Errorf("診断パッケージ用公開鍵は暗号化用途に対応していません。暗号化用の公開鍵を同梱する必要があります")
	}

	zipData, manifest, err := buildDiagnosticZip(configPath, cfg)
	if err != nil {
		return "", err
	}

	encrypted, err := encryptDiagnosticZip(zipData, entities)
	if err != nil {
		return "", err
	}

	outputDir := filepath.Dir(configPath)
	if strings.TrimSpace(configPath) == "" {
		outputDir = "."
	}
	outputPath := filepath.Join(outputDir, fmt.Sprintf("ClipForVRChat-diagnostics-%s.zip.gpg", time.Now().Format("20060102-150405")))
	if err := appcore.WritePrivateFile(outputPath, encrypted); err != nil {
		return "", fmt.Errorf("診断パッケージを保存できません: %w", err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "diagnostic package=%q files=%d missing=%d", outputPath, len(manifest.Files), len(manifest.MissingFiles))
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
	exePath, err := os.Executable()
	if err != nil {
		exePath = ""
	}
	exeSHA256 := ""
	if exePath != "" {
		exeSHA256, _ = fileSHA256(exePath)
	}
	configJSON := json.RawMessage(configSummaryForLog(cfg))
	manifest := diagnosticPackageManifest{
		CreatedAt:   time.Now().Format(time.RFC3339),
		AppVersion:  appVersion(),
		Version:     version,
		Revision:    revision,
		ReleaseTime: appReleaseTime(),
		Executable:  exePath,
		ExeSHA256:   exeSHA256,
		ConfigPath:  configPath,
		HistoryPath: appcore.HistoryPath(configPath),
		LogPath:     appcore.DiagnosticLogPath(configPath),
		Config:      configJSON,
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
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
		manifest.Files = append(manifest.Files, diagnosticPackageFile{Name: name, Path: path, SHA256: hash, Size: info.Size()})
		return addFileToZip(zipWriter, name, path)
	}

	if err := add("config.json", configPath); err != nil {
		_ = zipWriter.Close()
		return nil, manifest, err
	}
	if err := add("history.json", appcore.HistoryPath(configPath)); err != nil {
		_ = zipWriter.Close()
		return nil, manifest, err
	}
	if exePath != "" {
		if err := add("ClipForVRChat.exe", exePath); err != nil {
			_ = zipWriter.Close()
			return nil, manifest, err
		}
	} else {
		manifest.MissingFiles = append(manifest.MissingFiles, "ClipForVRChat.exe")
	}
	logFiles, _ := filepath.Glob(filepath.Join(appcore.DiagnosticLogDir(configPath), "*.log"))
	for _, logPath := range logFiles {
		name := filepath.Join("logs", filepath.Base(logPath))
		if err := add(name, logPath); err != nil {
			_ = zipWriter.Close()
			return nil, manifest, err
		}
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		_ = zipWriter.Close()
		return nil, manifest, err
	}
	manifestHeader := &zip.FileHeader{Name: "manifest.json", Method: zip.Deflate}
	manifestHeader.SetMode(0600)
	manifestWriter, err := zipWriter.CreateHeader(manifestHeader)
	if err != nil {
		_ = zipWriter.Close()
		return nil, manifest, err
	}
	if _, err := manifestWriter.Write(append(manifestData, '\n')); err != nil {
		_ = zipWriter.Close()
		return nil, manifest, err
	}
	if err := zipWriter.Close(); err != nil {
		return nil, manifest, err
	}
	return buf.Bytes(), manifest, nil
}

func addFileToZip(zipWriter *zip.Writer, name, path string) error {
	file, err := os.Open(path) // #nosec G304 -- diagnostic package paths are derived from app-owned config/history/log paths and the running executable.
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
