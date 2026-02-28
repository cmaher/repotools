package metrics

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountLOC_NoMarker(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("line1\nline2\nline3\n"), 0644)

	n, err := CountLOC(path, "")
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("got %d, want 3", n)
	}
}

func TestCountLOC_WithMarker(t *testing.T) {
	n, err := CountLOC("../../testdata/fixtures/sample.rs", `^#\[cfg\(test\)\]`)
	if err != nil {
		t.Fatal(err)
	}
	// Lines before #[cfg(test)] = 10 (including blank line before it)
	if n != 10 {
		t.Errorf("got %d, want 10", n)
	}
}

func TestCountLOC_AutoDetectRust(t *testing.T) {
	n, err := CountLOCAutoDetect("../../testdata/fixtures/sample.rs", "")
	if err != nil {
		t.Fatal(err)
	}
	if n != 10 {
		t.Errorf("got %d, want 10", n)
	}
}

func TestResolveFiles_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	os.WriteFile(path, []byte(""), 0644)

	files, err := ResolveFiles([]string{path}, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != path {
		t.Errorf("got %v", files)
	}
}

func TestResolveFiles_Directory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "b.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte(""), 0644)

	files, err := ResolveFiles([]string{dir}, "*.go", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("got %d files, want 2: %v", len(files), files)
	}
}

func TestResolveFiles_Exclude(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "a_test.go"), []byte(""), 0644)

	files, err := ResolveFiles([]string{dir}, "*.go", "_test\\.go$")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("got %d files, want 1: %v", len(files), files)
	}
}

func TestRunLOC_Output(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("one\ntwo\n"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("one\ntwo\nthree\n"), 0644)

	var buf bytes.Buffer
	err := RunLOC(&buf, []string{filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt")}, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 ") {
		t.Errorf("missing count for a.txt in:\n%s", out)
	}
	if !strings.Contains(out, "3 ") {
		t.Errorf("missing count for b.txt in:\n%s", out)
	}
	if !strings.Contains(out, "5 total") {
		t.Errorf("missing total in:\n%s", out)
	}
}
