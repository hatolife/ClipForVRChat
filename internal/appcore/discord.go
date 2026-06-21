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
	Attachments []struct {
		URL string `json:"url"`
	} `json:"attachments"`
}

func UploadDiscord(webhookURL string, filename string, encoded EncodedImage) (string, error) {
	if strings.TrimSpace(webhookURL) == "" {
		return "", errors.New("Discord投稿がONですがWebhook URLが未設定です。設定画面でWebhook URLを設定してください。")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files[0]", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(encoded.Data); err != nil {
		return "", err
	}
	if err := writer.WriteField("payload_json", `{"content":""}`); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL+"?wait=true", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Discord投稿に失敗しました。Webhook URLと投稿権限を確認してください。status=%d body=%s", resp.StatusCode, string(respData))
	}

	var parsed webhookResponse
	if err := json.Unmarshal(respData, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Attachments) == 0 || parsed.Attachments[0].URL == "" {
		return "", errors.New("Discord投稿は成功しましたが、画像URLを取得できませんでした。")
	}
	return parsed.Attachments[0].URL, nil
}
