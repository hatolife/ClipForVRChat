package appcore

import "testing"

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
