package appcore

import (
	"errors"
	"image"
	"net/url"
	"strings"

	"github.com/liyue201/goqr"
)

func DetectQRCodeURLs(img image.Image) []string {
	codes, err := goqr.Recognize(img)
	if err != nil {
		if errors.Is(err, goqr.ErrNoQRCode) {
			return nil
		}
		return nil
	}
	seen := map[string]bool{}
	var urls []string
	for _, code := range codes {
		raw := strings.TrimSpace(string(code.Payload))
		if !isHTTPURL(raw) || seen[raw] {
			continue
		}
		seen[raw] = true
		urls = append(urls, raw)
	}
	return urls
}

func isHTTPURL(raw string) bool {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return false
	}
	return (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Hostname() != ""
}

func qrDiscordContent(urls []string) string {
	if len(urls) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("QRコードのURL:\n")
	for _, u := range urls {
		b.WriteString(u)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}
