package backup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "test.db")
	content := []byte("hello world")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	backupPath, err := Create(src)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	// Verify backup exists
	got, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("backup content = %q, want %q", got, content)
	}
}

func TestCreatePathFormat(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "test.db")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	backupPath, err := Create(src)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	// Path should be <src>.<timestamp>.bak
	if !strings.HasPrefix(backupPath, src+".") {
		t.Errorf("backup path %q should start with %q", backupPath, src+".")
	}
	if !strings.HasSuffix(backupPath, ".bak") {
		t.Errorf("backup path %q should end with .bak", backupPath)
	}
}

func TestCreateMissingSource(t *testing.T) {
	_, err := Create("/nonexistent/path/test.db")
	if err == nil {
		t.Fatal("Create() should fail for missing source file")
	}
}
