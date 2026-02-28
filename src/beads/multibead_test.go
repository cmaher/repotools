package beads

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatMultiBeadHeader(t *testing.T) {
	var buf bytes.Buffer
	FormatMultiBeadHeader(&buf, "ABC-123")
	out := buf.String()
	if !strings.Contains(out, "=== ABC-123 ===") {
		t.Errorf("missing header in: %q", out)
	}
}
