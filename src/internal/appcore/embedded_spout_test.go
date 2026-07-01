package appcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedSpoutHelperDirUsesCombinedHash(t *testing.T) {
	t.Cleanup(func() { spoutHelperCacheDirForTest = "" })
	spoutHelperCacheDirForTest = t.TempDir()

	got, err := embeddedSpoutHelperDir([]byte("helper-v1"), []byte("dll-v1"))
	if err != nil {
		t.Fatal(err)
	}
	again, err := embeddedSpoutHelperDir([]byte("helper-v1"), []byte("dll-v1"))
	if err != nil {
		t.Fatal(err)
	}
	changed, err := embeddedSpoutHelperDir([]byte("helper-v2"), []byte("dll-v1"))
	if err != nil {
		t.Fatal(err)
	}
	if got != again {
		t.Fatalf("dir = %q, want stable %q", again, got)
	}
	if got == changed {
		t.Fatalf("dir did not change when helper bytes changed: %q", got)
	}
	changedDLL, err := embeddedSpoutHelperDir([]byte("helper-v1"), []byte("dll-v2"))
	if err != nil {
		t.Fatal(err)
	}
	if got == changedDLL {
		t.Fatalf("dir did not change when dll bytes changed: %q", got)
	}
	if filepath.Dir(got) != spoutHelperCacheDirForTest {
		t.Fatalf("dir parent = %q, want %q", filepath.Dir(got), spoutHelperCacheDirForTest)
	}
}

func TestEnsureEmbeddedSpoutFileWritesAndRepairs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spout-capture.exe")
	data := []byte("helper-bytes")
	if err := ensureEmbeddedSpoutFile(path, data, true); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(data) {
		t.Fatalf("data = %q, want %q", got, data)
	}
	if err := os.WriteFile(path, []byte("tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := ensureEmbeddedSpoutFile(path, data, true); err != nil {
		t.Fatal(err)
	}
	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(data) {
		t.Fatalf("repaired data = %q, want %q", got, data)
	}
}

func TestEnsureEmbeddedSpoutFileRejectsEmptyAsset(t *testing.T) {
	err := ensureEmbeddedSpoutFile(filepath.Join(t.TempDir(), spoutHelperFileName), nil, true)
	if err == nil {
		t.Fatal("err = nil, want empty asset error")
	}
}
