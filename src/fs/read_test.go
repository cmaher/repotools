package fs

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadLines_Full(t *testing.T) {
	path := writeTempFile(t, "line1\nline2\nline3\n")
	var buf bytes.Buffer
	err := ReadLines(&buf, path, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	if !strings.Contains(lines[0], "1\tline1") {
		t.Errorf("line 0 = %q, want line number + tab + content", lines[0])
	}
}

func TestReadLines_Range(t *testing.T) {
	path := writeTempFile(t, "a\nb\nc\nd\ne\n")
	var buf bytes.Buffer
	err := ReadLines(&buf, path, 2, 4)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	if !strings.Contains(lines[0], "2\tb") {
		t.Errorf("first line = %q, want line 2", lines[0])
	}
	if !strings.Contains(lines[2], "4\td") {
		t.Errorf("last line = %q, want line 4", lines[2])
	}
}

func TestReadLines_NotFound(t *testing.T) {
	var buf bytes.Buffer
	err := ReadLines(&buf, "/nonexistent/file.txt", 0, 0)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadLines_IsDirectory(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	err := ReadLines(&buf, dir, 0, 0)
	if err == nil {
		t.Fatal("expected error for directory")
	}
}
