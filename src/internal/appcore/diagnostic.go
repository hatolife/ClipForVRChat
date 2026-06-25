package appcore

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	discordWebhookURLPattern = regexp.MustCompile(`https://(?:discord|discordapp)\.com/api/webhooks/[^\s"'<>\\]+`)
	jsonSecretFieldPattern   = regexp.MustCompile(`(?i)("(?:discordToken|webhookUrl)"\s*:\s*")[^"]*(")`)
)

func DiagnosticLogPath(configPath string) string {
	name := time.Now().Format("2006-01-02") + ".log"
	return filepath.Join(DiagnosticLogDir(configPath), name)
}

func DiagnosticLogDir(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return "logs"
	}
	return filepath.Join(filepath.Dir(configPath), "logs")
}

func AppendDiagnosticLog(path string, format string, args ...any) {
	if strings.TrimSpace(path) == "" {
		return
	}
	line := RedactDiagnosticText(fmt.Sprintf(format, args...))
	entry := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), line)
	if err := os.MkdirAll(filepath.Dir(path), privateDirMode); err != nil {
		return
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, privateFileMode) // #nosec G304 -- diagnostic log path is derived from the active config path and uses private permissions.
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(entry)
	_ = file.Chmod(privateFileMode)
}

func RedactDiagnosticText(text string) string {
	if text == "" {
		return text
	}
	text = redactKnownUserPaths(text)
	text = discordWebhookURLPattern.ReplaceAllString(text, "<redacted-discord-webhook-url>")
	text = jsonSecretFieldPattern.ReplaceAllString(text, `${1}<redacted>${2}`)
	return text
}

func redactKnownUserPaths(text string) string {
	candidates := []struct {
		name  string
		value string
	}{
		{"USERPROFILE", os.Getenv("USERPROFILE")},
		{"APPDATA", os.Getenv("APPDATA")},
		{"LOCALAPPDATA", os.Getenv("LOCALAPPDATA")},
		{"TEMP", os.Getenv("TEMP")},
		{"TMP", os.Getenv("TMP")},
		{"HOME", os.Getenv("HOME")},
	}
	replacements := make([]string, 0, len(candidates)*8)
	seen := map[string]bool{}
	for _, candidate := range candidates {
		value := strings.Trim(strings.TrimSpace(candidate.value), `"`)
		if value == "" {
			continue
		}
		placeholder := "%" + candidate.name + "%"
		variants := []string{filepath.Clean(value)}
		slash := strings.ReplaceAll(filepath.Clean(value), `\`, `/`)
		backslash := strings.ReplaceAll(filepath.Clean(value), `/`, `\`)
		escapedBackslash := strings.ReplaceAll(backslash, `\`, `\\`)
		variants = append(variants, slash, backslash, escapedBackslash)
		for _, variant := range variants {
			if variant == "" || variant == "." || seen[variant] {
				continue
			}
			seen[variant] = true
			replacements = append(replacements, variant, placeholder)
		}
	}
	if len(replacements) == 0 {
		return text
	}
	return strings.NewReplacer(replacements...).Replace(text)
}
