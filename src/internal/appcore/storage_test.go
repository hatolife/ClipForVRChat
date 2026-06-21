package appcore

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSaveConfigUsesPrivateFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX file mode semantics")
	}
	requireEffectiveFileModes(t)
	path := filepath.Join(t.TempDir(), "config.json")
	if err := SaveConfig(path, DefaultConfig()); err != nil {
		t.Fatal(err)
	}
	assertFileMode(t, path, privateFileMode)
}

func TestSaveHistoryUsesPrivateFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX file mode semantics")
	}
	requireEffectiveFileModes(t)
	path := filepath.Join(t.TempDir(), "history.json")
	if err := SaveHistory(path, []HistoryEntry{{ID: "1"}}); err != nil {
		t.Fatal(err)
	}
	assertFileMode(t, path, privateFileMode)
}

func requireEffectiveFileModes(t *testing.T) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mode-check")
	if err := os.WriteFile(path, []byte("x"), privateFileMode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, privateFileMode); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != privateFileMode {
		t.Skip("filesystem does not preserve POSIX file mode semantics")
	}
}

func assertFileMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("mode = %v, want %v", got, want)
	}
}
