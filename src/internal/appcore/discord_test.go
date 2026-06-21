package appcore

import (
	"strings"
	"testing"
)

func TestValidateDiscordWebhookURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ok   bool
	}{
		{name: "discord", raw: "https://discord.com/api/webhooks/123/abc", ok: true},
		{name: "discordapp", raw: "https://discordapp.com/api/webhooks/123/abc", ok: true},
		{name: "http", raw: "http://discord.com/api/webhooks/123/abc", ok: false},
		{name: "wrong host", raw: "https://example.com/api/webhooks/123/abc", ok: false},
		{name: "wrong path", raw: "https://discord.com/webhooks/123/abc", ok: false},
		{name: "missing token", raw: "https://discord.com/api/webhooks/123", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, postURL, err := ValidateDiscordWebhookURL(tt.raw)
			if tt.ok && err != nil {
				t.Fatalf("expected valid URL: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("expected invalid URL")
			}
			if tt.ok && postURL == tt.raw {
				t.Fatal("expected wait=true to be added to post URL")
			}
		})
	}
}

func TestIsTrustedDiscordImageURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ok   bool
	}{
		{name: "cdn attachment", raw: "https://cdn.discordapp.com/attachments/1/2/image.png", ok: true},
		{name: "media attachment", raw: "https://media.discordapp.net/attachments/1/2/image.png", ok: true},
		{name: "http", raw: "http://cdn.discordapp.com/attachments/1/2/image.png", ok: false},
		{name: "wrong host", raw: "https://example.com/attachments/1/2/image.png", ok: false},
		{name: "wrong path", raw: "https://cdn.discordapp.com/icons/1/image.png", ok: false},
		{name: "local", raw: "http://127.0.0.1/attachments/1/2/image.png", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTrustedDiscordImageURL(tt.raw); got != tt.ok {
				t.Fatalf("IsTrustedDiscordImageURL() = %v, want %v", got, tt.ok)
			}
		})
	}
}

func TestParseWebhookURL(t *testing.T) {
	id, token := ParseWebhookURL("https://discord.com/api/webhooks/123/abc")
	if id != "123" || token != "abc" {
		t.Fatalf("ParseWebhookURL = %q, %q", id, token)
	}

	id, token = ParseWebhookURL("https://example.com/api/webhooks/123/abc")
	if id != "" || token != "" {
		t.Fatalf("invalid ParseWebhookURL = %q, %q", id, token)
	}
}

func TestDeleteDiscordMessageRequiresFields(t *testing.T) {
	if err := DeleteDiscordMessage("", "token", "message"); err == nil {
		t.Fatal("expected missing webhook ID error")
	}
	if err := DeleteDiscordMessage("id", "", "message"); err == nil {
		t.Fatal("expected missing token error")
	}
	if err := DeleteDiscordMessage("id", "token", ""); err == nil {
		t.Fatal("expected missing message ID error")
	}
}

func TestDiscordUploadStatusErrorGuidesWebhookRefresh(t *testing.T) {
	err := discordUploadStatusError(401, []byte(`{"message": "Invalid Webhook Token", "code": 50027}`))
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Webhook URLが無効") {
		t.Fatalf("message should explain invalid webhook: %q", msg)
	}
	if !strings.Contains(msg, "再度コピー") || !strings.Contains(msg, "貼り直") {
		t.Fatalf("message should guide re-copying URL: %q", msg)
	}
	if strings.Contains(msg, "body=") {
		t.Fatalf("message should not expose raw body for invalid token: %q", msg)
	}
}

func TestDiscordUploadStatusErrorTruncatesUnexpectedBody(t *testing.T) {
	err := discordUploadStatusError(500, []byte(strings.Repeat("x", 200)))
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "body=") {
		t.Fatalf("message should keep limited body: %q", msg)
	}
	if len([]rune(msg)) > 260 {
		t.Fatalf("message too long: %d runes", len([]rune(msg)))
	}
}

func TestDiscordWebhookSetupMessage(t *testing.T) {
	msg := discordWebhookSetupMessage("Discord Webhook URLが未設定です。")
	if !strings.Contains(msg, "再度コピー") || !strings.Contains(msg, "自動投稿用Webhook URL") {
		t.Fatalf("unexpected setup message: %q", msg)
	}
}
