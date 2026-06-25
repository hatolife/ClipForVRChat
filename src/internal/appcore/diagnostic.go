package appcore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	line := fmt.Sprintf(format, args...)
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
