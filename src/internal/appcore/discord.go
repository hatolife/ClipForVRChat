package appcore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type webhookResponse struct {
	ID          string `json:"id"`
	WebhookID   string `json:"webhook_id"`
	Attachments []struct {
		URL string `json:"url"`
	} `json:"attachments"`
}

type DiscordUpload struct {
	URL       string
	MessageID string
	WebhookID string
	Token     string
}

type discordAllowedMentions struct {
	Parse []string `json:"parse"`
}

type discordWebhookPayload struct {
	Content         string                 `json:"content"`
	AllowedMentions discordAllowedMentions `json:"allowed_mentions"`
}

func UploadDiscord(webhookURL string, filename string, encoded EncodedImage) (DiscordUpload, error) {
	return UploadDiscordWithContent(webhookURL, filename, encoded, "")
}

func UploadDiscordWithContent(webhookURL string, filename string, encoded EncodedImage, content string) (DiscordUpload, error) {
	var uploaded DiscordUpload
	if strings.TrimSpace(webhookURL) == "" {
		return uploaded, errors.New(discordWebhookSetupMessage("Discord Webhook URLが未設定です。"))
	}
	webhookID, token, postURL, err := ValidateDiscordWebhookURL(webhookURL)
	if err != nil {
		return uploaded, errors.New(discordWebhookSetupMessage(err.Error()))
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files[0]", filename)
	if err != nil {
		return uploaded, err
	}
	if _, err := part.Write(encoded.Data); err != nil {
		return uploaded, err
	}
	payload, err := discordPayloadJSON(content)
	if err != nil {
		return uploaded, err
	}
	if err := writer.WriteField("payload_json", string(payload)); err != nil {
		return uploaded, err
	}
	if err := writer.Close(); err != nil {
		return uploaded, err
	}

	req, err := http.NewRequest(http.MethodPost, postURL, &body)
	if err != nil {
		return uploaded, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return uploaded, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return uploaded, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return uploaded, discordUploadStatusError(resp.StatusCode, respData)
	}

	var parsed webhookResponse
	if err := json.Unmarshal(respData, &parsed); err != nil {
		return uploaded, err
	}
	if len(parsed.Attachments) == 0 || parsed.Attachments[0].URL == "" {
		return uploaded, errors.New("Discord投稿は成功しましたが、画像URLを取得できませんでした。")
	}
	if !IsTrustedDiscordImageURL(parsed.Attachments[0].URL) {
		return uploaded, errors.New("Discord投稿は成功しましたが、取得した画像URLの形式を確認できませんでした。")
	}
	uploaded.URL = parsed.Attachments[0].URL
	uploaded.MessageID = parsed.ID
	uploaded.WebhookID = parsed.WebhookID
	if uploaded.WebhookID == "" {
		uploaded.WebhookID = webhookID
	}
	uploaded.Token = token
	return uploaded, nil
}

func discordPayloadJSON(content string) ([]byte, error) {
	return json.Marshal(discordWebhookPayload{
		Content: truncateDiscordContent(content),
		AllowedMentions: discordAllowedMentions{
			Parse: []string{},
		},
	})
}

func truncateDiscordContent(content string) string {
	const maxDiscordContentLength = 2000
	runes := []rune(strings.TrimSpace(content))
	if len(runes) <= maxDiscordContentLength {
		return string(runes)
	}
	return string(runes[:maxDiscordContentLength-3]) + "..."
}

func discordUploadStatusError(status int, body []byte) error {
	bodyText := strings.TrimSpace(string(body))
	if status == http.StatusUnauthorized || status == http.StatusNotFound || strings.Contains(bodyText, "Invalid Webhook Token") {
		return fmt.Errorf("%s status=%d", discordWebhookSetupMessage("Discord Webhook URLが無効になっている可能性があります。"), status)
	}
	if status == http.StatusForbidden {
		return fmt.Errorf("Discord投稿に失敗しました。Webhookの投稿権限を確認してください。Webhook URLをDiscordで再取得して設定し直すと改善する場合があります。status=%d", status)
	}
	if bodyText != "" {
		return fmt.Errorf("Discord投稿に失敗しました。時間をおいて再試行してください。改善しない場合はWebhook URLをDiscordで再取得して設定し直してください。status=%d body=%s", status, truncateForUserMessage(bodyText, 160))
	}
	return fmt.Errorf("Discord投稿に失敗しました。時間をおいて再試行してください。改善しない場合はWebhook URLをDiscordで再取得して設定し直してください。status=%d", status)
}

func discordWebhookSetupMessage(reason string) string {
	return reason + " DiscordでWebhook URLを再度コピーし、設定画面の「Discord Webhook URL」または「自動投稿用Webhook URL」に貼り直してください。"
}

func truncateForUserMessage(text string, maxRunes int) string {
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes]) + "..."
}

func DeleteDiscordMessage(webhookID, token, messageID string) error {
	if webhookID == "" || token == "" || messageID == "" {
		return errors.New("Discord削除に必要なWebhook情報が履歴にありません。")
	}
	url := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s/messages/%s", webhookID, token, messageID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil
	}
	respData, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("Discord画像を削除できませんでした。status=%d body=%s", resp.StatusCode, string(respData))
}

func ParseWebhookURL(webhookURL string) (string, string) {
	webhookID, token, _, err := ValidateDiscordWebhookURL(webhookURL)
	if err != nil {
		return "", ""
	}
	return webhookID, token
}

func ValidateDiscordWebhookURL(webhookURL string) (string, string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(webhookURL))
	if err != nil {
		return "", "", "", errors.New("Discord Webhook URLの形式が正しくありません。")
	}
	if parsed.Scheme != "https" {
		return "", "", "", errors.New("Discord Webhook URLは https:// で始まるURLを指定してください。")
	}
	host := strings.ToLower(parsed.Hostname())
	if host != "discord.com" && host != "discordapp.com" {
		return "", "", "", errors.New("Discord Webhook URLは discord.com または discordapp.com のURLを指定してください。")
	}
	parts := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")
	if len(parts) != 4 || parts[0] != "api" || parts[1] != "webhooks" {
		return "", "", "", errors.New("Discord Webhook URLは /api/webhooks/{id}/{token} 形式を指定してください。")
	}
	webhookID, err := url.PathUnescape(parts[2])
	if err != nil || strings.TrimSpace(webhookID) == "" {
		return "", "", "", errors.New("Discord Webhook URLのWebhook IDを確認できません。")
	}
	token, err := url.PathUnescape(parts[3])
	if err != nil || strings.TrimSpace(token) == "" {
		return "", "", "", errors.New("Discord Webhook URLのtokenを確認できません。")
	}
	query := parsed.Query()
	query.Set("wait", "true")
	parsed.RawQuery = query.Encode()
	return webhookID, token, parsed.String(), nil
}

func IsTrustedDiscordImageURL(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme != "https" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	switch host {
	case "cdn.discordapp.com", "media.discordapp.net":
	default:
		return false
	}
	parts := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")
	if len(parts) < 3 || parts[0] != "attachments" {
		return false
	}
	return true
}
