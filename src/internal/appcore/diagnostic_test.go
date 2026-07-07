package appcore

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRedactDiagnosticTextRedactsUserPathsAndSecrets(t *testing.T) {
	userProfile := filepath.Join(t.TempDir(), "Users", "alice")
	t.Setenv("USERPROFILE", userProfile)
	text := `path="` + filepath.Join(userProfile, "Pictures", "image.png") + `" webhook="https://discord.com/api/webhooks/123/token" json={"discordToken":"secret","webhookUrl":"https://discordapp.com/api/webhooks/456/token"}`

	got := RedactDiagnosticText(text)

	for _, leaked := range []string{userProfile, "https://discord.com/api/webhooks/123/token", "https://discordapp.com/api/webhooks/456/token", `"discordToken":"secret"`} {
		if strings.Contains(got, leaked) {
			t.Fatalf("RedactDiagnosticText leaked %q in %q", leaked, got)
		}
	}
	if !strings.Contains(got, "%USERPROFILE%") || !strings.Contains(got, "<redacted>") {
		t.Fatalf("RedactDiagnosticText = %q, want placeholders", got)
	}
}

func TestRedactDiagnosticTextRedactsWebhookURLInsideErrors(t *testing.T) {
	text := `Post "https://discord.com/api/webhooks/1234567890/live-token": dial tcp: lookup discord.com: no such host`

	got := RedactDiagnosticText(text)

	if strings.Contains(got, "live-token") || strings.Contains(got, "1234567890") {
		t.Fatalf("RedactDiagnosticText leaked webhook URL in %q", got)
	}
	if !strings.Contains(got, "<redacted-discord-webhook-url>") {
		t.Fatalf("RedactDiagnosticText = %q, want webhook placeholder", got)
	}
}
