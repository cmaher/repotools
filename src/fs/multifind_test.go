package fs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultiFind(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		os.WriteFile(filepath.Join(dir, name), []byte(""), 0644)
	}

	var buf bytes.Buffer
	MultiFind(&buf, 10, []string{"-name", "*.txt"}, []string{dir})
	out := buf.String()

	if !strings.Contains(out, "==> "+dir+" <==") {
		t.Errorf("missing header")
	}
	if !strings.Contains(out, "a.txt") {
		t.Errorf("missing a.txt")
	}
}

func TestMultiFind_HeadLimit(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 20; i++ {
		os.WriteFile(filepath.Join(dir, fmt.Sprintf("file%02d.txt", i)), []byte(""), 0644)
	}

	var buf bytes.Buffer
	MultiFind(&buf, 3, []string{"-name", "*.txt"}, []string{dir})
	out := buf.String()

	lines := strings.Split(strings.TrimSpace(out), "\n")
	contentLines := 0
	for _, l := range lines {
		if l != "" && !strings.HasPrefix(l, "==>") && l != "---" {
			contentLines++
		}
	}
	if contentLines > 3 {
		t.Errorf("got %d content lines, want <= 3", contentLines)
	}
}
