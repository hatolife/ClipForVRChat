package appcore

import (
	"encoding/json"
	"os"
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

func TestAddResultsToHistoryKeepsLocalAndQRCodeResults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	results := []Result{
		{Name: "local.png", OutputPath: filepath.Join(t.TempDir(), "local.png")},
		{Name: "qr.png", QRURLs: []string{"https://example.com/qr"}},
		{Name: "empty.png"},
	}

	history, err := AddResultsToHistory(path, results)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(history))
	}
	if history[0].Name != "qr.png" || len(history[0].QRURLs) != 1 {
		t.Fatalf("first history = %+v, want qr entry", history[0])
	}
	if history[1].Name != "local.png" || history[1].OutputPath == "" {
		t.Fatalf("second history = %+v, want local entry", history[1])
	}
	if results[2].HistoryID != "" {
		t.Fatal("empty result should not receive history ID")
	}
}

func TestRemoteURLAvailableRejectsUntrustedURLWithoutRequest(t *testing.T) {
	if RemoteURLAvailable("http://127.0.0.1:1/attachments/1/2/image.png") {
		t.Fatal("untrusted URL should not be available")
	}
}

func TestHistoryPathUsesConfigDirectory(t *testing.T) {
	got := HistoryPath(filepath.Join("base", "config.json"))
	if filepath.Base(got) != "history.json" {
		t.Fatalf("HistoryPath base = %q, want history.json", filepath.Base(got))
	}
	if HistoryPath("") != "history.json" {
		t.Fatalf("HistoryPath empty = %q, want history.json", HistoryPath(""))
	}
}

func TestLoadHistoryMissingAndInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	missing, err := LoadHistory(filepath.Join(dir, "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if missing != nil {
		t.Fatalf("missing history = %+v, want nil", missing)
	}

	invalid := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalid, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadHistory(invalid); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestMarkHistoryCleared(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"},
		{ID: "2", URL: "https://cdn.discordapp.com/attachments/1/2/b.png"},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, err := MarkHistoryCleared(path, []string{"2", "missing"})
	if err != nil {
		t.Fatal(err)
	}
	if history[0].Cleared {
		t.Fatal("entry 1 should not be cleared")
	}
	if !history[1].Cleared || history[1].ClearedAt == "" {
		t.Fatalf("entry 2 should be cleared with timestamp: %+v", history[1])
	}
}

func TestSetHistoryPinned(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, err := SetHistoryPinned(path, "1", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || !history[0].Pinned {
		t.Fatalf("history = %+v, want pinned entry", history)
	}
}

func TestDeleteHistoryEntriesKeepsPinnedEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "delete", Name: "delete.png"},
		{ID: "pinned", Name: "pinned.png", Pinned: true},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, removed, err := DeleteHistoryEntries(path, []string{"delete", "pinned"})
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 || len(history) != 1 || history[0].ID != "pinned" {
		t.Fatalf("history = %+v removed = %d, want only pinned kept", history, removed)
	}
}

func TestDeleteLocalHistoryFilesMarksDeleted(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.png")
	if err := os.WriteFile(output, []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "history.json")
	initial := []HistoryEntry{
		{ID: "delete", OutputPath: output},
		{ID: "pinned", OutputPath: output, Pinned: true},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, deleted, err := DeleteLocalHistoryFiles(path, []string{"delete", "pinned"})
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 1 || !history[0].LocalDeleted || history[0].LocalDeletedAt == "" || history[0].LocalExists {
		t.Fatalf("history = %+v deleted = %d, want first local deleted", history, deleted)
	}
	if history[1].LocalDeleted {
		t.Fatalf("pinned entry should not be deleted: %+v", history[1])
	}
}

func TestPurgeUnavailableHistoryRemovesUntrustedURLs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "untrusted", URL: "https://example.com/attachments/1/2/b.png"},
		{ID: "empty", URL: ""},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, removed, err := PurgeUnavailableHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("removed = %d, want 1", removed)
	}
	if len(history) != 1 || history[0].ID != "empty" {
		t.Fatalf("history = %+v, want only empty entry", history)
	}
}

func TestPurgeUnavailableHistoryKeepsPinnedEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "pinned", URL: "https://example.com/attachments/1/2/b.png", Pinned: true},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, removed, err := PurgeUnavailableHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 0 || len(history) != 1 || history[0].ID != "pinned" {
		t.Fatalf("history = %+v removed = %d, want pinned entry kept", history, removed)
	}
}

func TestPurgeDiscordDeletedHistoryRemovesOnlyDeletedEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := []HistoryEntry{
		{ID: "deleted", URL: "https://cdn.discordapp.com/attachments/1/2/a.png", DiscordDeleted: true},
		{ID: "available", URL: "https://cdn.discordapp.com/attachments/1/2/b.png"},
		{ID: "pinned", URL: "https://cdn.discordapp.com/attachments/1/2/c.png", DiscordDeleted: true, Pinned: true},
	}
	if err := SaveHistory(path, initial); err != nil {
		t.Fatal(err)
	}

	history, removed, err := PurgeDiscordDeletedHistory(path, false)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("removed = %d, want 1", removed)
	}
	if len(history) != 2 || history[0].ID != "available" || history[1].ID != "pinned" {
		t.Fatalf("history = %+v, want available and pinned", history)
	}
}

func TestPurgeDiscordDeletedHistoryDeletesOutputWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.png")
	if err := os.WriteFile(output, []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "history.json")
	if err := SaveHistory(path, []HistoryEntry{{ID: "deleted", DiscordDeleted: true, OutputPath: output}}); err != nil {
		t.Fatal(err)
	}

	if _, _, err := PurgeDiscordDeletedHistory(path, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("output file should be removed, stat err = %v", err)
	}
}

func TestPurgeDiscordDeletedHistoryKeepsOutputWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.png")
	if err := os.WriteFile(output, []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "history.json")
	if err := SaveHistory(path, []HistoryEntry{{ID: "deleted", DiscordDeleted: true, OutputPath: output}}); err != nil {
		t.Fatal(err)
	}

	if _, _, err := PurgeDiscordDeletedHistory(path, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("output file should remain, stat err = %v", err)
	}
}

func TestSaveHistoryWritesJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := SaveHistory(path, []HistoryEntry{{ID: "1", Name: "image.png"}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded []HistoryEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 1 || decoded[0].ID != "1" {
		t.Fatalf("decoded = %+v", decoded)
	}
}
