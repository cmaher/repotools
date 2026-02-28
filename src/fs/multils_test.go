package fs

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultiLS(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	os.WriteFile(filepath.Join(dir1, "a.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir2, "b.txt"), []byte(""), 0644)

	var buf bytes.Buffer
	MultiLS(&buf, []string{dir1, dir2})
	out := buf.String()

	if !strings.Contains(out, "==> "+dir1+" <==") {
		t.Errorf("missing header for dir1")
	}
	if !strings.Contains(out, "a.txt") {
		t.Errorf("missing a.txt in output")
	}
	if !strings.Contains(out, "==> "+dir2+" <==") {
		t.Errorf("missing header for dir2")
	}
	if !strings.Contains(out, "b.txt") {
		t.Errorf("missing b.txt in output")
	}
}
