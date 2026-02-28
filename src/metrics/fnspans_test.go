package metrics

import (
	"bytes"
	"strings"
	"testing"
)

func TestExtractFnSpans_Python(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.py", `^\s*(?:async\s+)?def\s+(\w+)`, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 3 {
		t.Fatalf("got %d spans, want 3: %v", len(spans), spans)
	}
	if spans[0].Name != "hello" {
		t.Errorf("first fn = %q, want hello", spans[0].Name)
	}
	if spans[2].Name != "fetch" {
		t.Errorf("third fn = %q, want fetch", spans[2].Name)
	}
}

func TestExtractFnSpans_Go(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.go", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 3 {
		t.Fatalf("got %d spans, want 3: %v", len(spans), spans)
	}
	names := []string{spans[0].Name, spans[1].Name, spans[2].Name}
	expected := []string{"main", "helper", "Handle"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("span %d name = %q, want %q", i, names[i], want)
		}
	}
}

func TestExtractFnSpans_SpanBoundaries(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.go", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	// main starts at line 3, helper starts at line 7
	// so main spans 3-6 (4 lines)
	if spans[0].Start != 3 || spans[0].End != 6 {
		t.Errorf("main span = %d-%d, want 3-6", spans[0].Start, spans[0].End)
	}
}

func TestExtractFnSpans_AfterFilter(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.rs", "", `^#\[cfg\(test\)\]`, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 1 || spans[0].Name != "test_main" {
		t.Errorf("got %v, want [test_main]", spans)
	}
}

func TestExtractFnSpans_IncludeFilter(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.py", `^\s*(?:async\s+)?def\s+(\w+)`, "", "hello", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 1 || spans[0].Name != "hello" {
		t.Errorf("got %v, want [hello]", spans)
	}
}

func TestExtractFnSpans_ExcludeFilter(t *testing.T) {
	spans, err := ExtractFnSpans("../../testdata/fixtures/sample.py", `^\s*(?:async\s+)?def\s+(\w+)`, "", "", "fetch")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 2 {
		t.Errorf("got %d spans, want 2", len(spans))
	}
}

func TestRunFnSpans_SingleFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunFnSpans(&buf, []string{"../../testdata/fixtures/sample.go"}, "", "", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "main") || !strings.Contains(out, "helper") {
		t.Errorf("output missing functions:\n%s", out)
	}
}

func TestRunFnSpans_MultiFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunFnSpans(&buf, []string{"../../testdata/fixtures/sample.go", "../../testdata/fixtures/sample.py"}, "", "", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "==>") {
		t.Errorf("multi-file output missing headers:\n%s", out)
	}
}
