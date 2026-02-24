package backup

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Create copies srcPath to srcPath.YYYYMMDD-HHMMSS.bak.
// Returns the backup file path.
func Create(srcPath string) (string, error) {
	ts := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.%s.bak", srcPath, ts)

	src, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("open source for backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(backupPath)
		return "", fmt.Errorf("backup copy failed: %w", err)
	}

	return backupPath, nil
}
