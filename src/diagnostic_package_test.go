package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncryptPathWithPublicKeyRejectsOversizedFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "oversized.bin")
	if err := os.WriteFile(filePath, []byte("01234567"), 0600); err != nil {
		t.Fatal(err)
	}
	oldMax := cliEncryptionMaxFileSize
	cliEncryptionMaxFileSize = 4
	t.Cleanup(func() {
		cliEncryptionMaxFileSize = oldMax
	})

	if _, err := encryptPathWithPublicKey(filePath); err == nil || !strings.Contains(err.Error(), "暗号化対象ファイルが大きすぎます") {
		t.Fatalf("expected oversize file error, got %v", err)
	}
}

func TestEncryptPathWithPublicKeyRejectsDirectoryFileCountLimit(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "report")
	if err := os.Mkdir(inputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "a.txt"), []byte("a"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "b.txt"), []byte("b"), 0600); err != nil {
		t.Fatal(err)
	}

	oldMax := cliEncryptionMaxDirectoryFiles
	cliEncryptionMaxDirectoryFiles = 1
	t.Cleanup(func() {
		cliEncryptionMaxDirectoryFiles = oldMax
	})

	if _, err := encryptPathWithPublicKey(inputDir); err == nil || !strings.Contains(err.Error(), "ファイル数が多すぎます") {
		t.Fatalf("expected file-count limit error, got %v", err)
	}
}

func TestEncryptPathWithPublicKeyRejectsDirectoryTotalSizeLimit(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "archive")
	if err := os.Mkdir(inputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "a.txt"), []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "b.txt"), []byte("def"), 0600); err != nil {
		t.Fatal(err)
	}

	oldMax := cliEncryptionMaxDirectorySize
	cliEncryptionMaxDirectorySize = 4
	t.Cleanup(func() {
		cliEncryptionMaxDirectorySize = oldMax
	})

	if _, err := encryptPathWithPublicKey(inputDir); err == nil || !strings.Contains(err.Error(), "総サイズが大きすぎます") {
		t.Fatalf("expected total-size limit error, got %v", err)
	}
}

func TestEncryptPathWithPublicKeyRejectsDirectoryDepthLimit(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "deep")
	if err := os.MkdirAll(filepath.Join(inputDir, "a", "b"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "a", "b", "file.txt"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	oldMax := cliEncryptionMaxDirectoryDepth
	cliEncryptionMaxDirectoryDepth = 2
	t.Cleanup(func() {
		cliEncryptionMaxDirectoryDepth = oldMax
	})

	if _, err := encryptPathWithPublicKey(inputDir); err == nil || !strings.Contains(err.Error(), "暗号化対象のフォルダが深すぎます") {
		t.Fatalf("expected directory-depth limit error, got %v", err)
	}
}
