package appcore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

func UploadDiscord(webhookURL string, filename string, encoded EncodedImage) (DiscordUpload, error) {
	var uploaded DiscordUpload
	if strings.TrimSpace(webhookURL) == "" {
		return uploaded, errors.New("Discord投稿がONですがWebhook URLが未設定です。設定画面でWebhook URLを設定してください。")
	}
	webhookID, token := ParseWebhookURL(webhookURL)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files[0]", filename)
	if err != nil {
		return uploaded, err
	}
	if _, err := part.Write(encoded.Data); err != nil {
		return uploaded, err
	}
	if err := writer.WriteField("payload_json", `{"content":""}`); err != nil {
		return uploaded, err
	}
	if err := writer.Close(); err != nil {
		return uploaded, err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL+"?wait=true", &body)
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
		return uploaded, fmt.Errorf("Discord投稿に失敗しました。Webhook URLと投稿権限を確認してください。status=%d body=%s", resp.StatusCode, string(respData))
	}

	var parsed webhookResponse
	if err := json.Unmarshal(respData, &parsed); err != nil {
		return uploaded, err
	}
	if len(parsed.Attachments) == 0 || parsed.Attachments[0].URL == "" {
		return uploaded, errors.New("Discord投稿は成功しましたが、画像URLを取得できませんでした。")
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
	parts := strings.Split(strings.TrimSpace(webhookURL), "/")
	for i := 0; i < len(parts)-2; i++ {
		if parts[i] == "webhooks" {
			return parts[i+1], strings.Split(parts[i+2], "?")[0]
		}
	}
	return "", ""
}
