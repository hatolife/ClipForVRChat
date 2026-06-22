package appcore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DiagnosticLogPath(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return "diagnostic.log"
	}
	return filepath.Join(filepath.Dir(configPath), "diagnostic.log")
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
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, privateFileMode)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(entry)
	_ = file.Chmod(privateFileMode)
}
