package appcore

import (
	"path/filepath"
	"testing"
)

func TestAddResultsToHistorySkipsUntrustedURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	results := []Result{
		{Name: "ok.png", URL: "https://cdn.discordapp.com/attachments/1/2/ok.png"},
		{Name: "bad.png", URL: "https://example.com/attachments/1/2/bad.png"},
	}

	history, err := AddResultsToHistory(path, results)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 {
		t.Fatalf("len(history) = %d, want 1", len(history))
	}
	if history[0].Name != "ok.png" {
		t.Fatalf("history[0].Name = %q, want ok.png", history[0].Name)
	}
	if results[0].HistoryID == "" {
		t.Fatal("trusted result should receive history ID")
	}
	if results[1].HistoryID != "" {
		t.Fatal("untrusted result should not receive history ID")
	}
}

func TestRemoteURLAvailableRejectsUntrustedURLWithoutRequest(t *testing.T) {
	if RemoteURLAvailable("http://127.0.0.1:1/attachments/1/2/image.png") {
		t.Fatal("untrusted URL should not be available")
	}
}
