//go:build windows && embeddedspout

package appcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveEmbeddedSpoutHelperWithEmbeddedAssets(t *testing.T) {
	t.Cleanup(func() { spoutHelperCacheDirForTest = "" })
	spoutHelperCacheDirForTest = t.TempDir()

	helper, err := resolveEmbeddedSpoutHelper()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(helper) != spoutHelperFileName {
		t.Fatalf("helper = %q, want %s", helper, spoutHelperFileName)
	}
	if _, err := os.Stat(helper); err != nil {
		t.Fatalf("helper was not extracted: %v", err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(helper), spoutLibraryDLLName)); err != nil {
		t.Fatalf("dll was not extracted next to helper: %v", err)
	}
}
