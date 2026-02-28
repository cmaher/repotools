package github

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestParsePRData(t *testing.T) {
	data, err := os.ReadFile("../../testdata/fixtures/pr.json")
	if err != nil {
		t.Fatal(err)
	}
	var pr PRData
	if err := json.Unmarshal(data, &pr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pr.Number != 42 {
		t.Errorf("number = %d, want 42", pr.Number)
	}
	if pr.Author.Login != "alice" {
		t.Errorf("author = %q, want alice", pr.Author.Login)
	}
	if len(pr.Comments) != 1 {
		t.Errorf("comments = %d, want 1", len(pr.Comments))
	}
}

func TestFilterSections_Only(t *testing.T) {
	got := FilterSections(AllSections, "info,body", "")
	if len(got) != 2 || got[0] != "info" || got[1] != "body" {
		t.Errorf("got %v", got)
	}
}

func TestFilterSections_Exclude(t *testing.T) {
	got := FilterSections(AllSections, "", "body,checks")
	for _, s := range got {
		if s == "body" || s == "checks" {
			t.Errorf("should not contain %q", s)
		}
	}
}

func TestFilterSections_InvalidOnly(t *testing.T) {
	_, err := ValidateSections("info,bogus")
	if err == nil {
		t.Fatal("expected error for invalid section")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error should mention 'bogus': %v", err)
	}
}

func TestRenderPR_AllSections(t *testing.T) {
	data, _ := os.ReadFile("../../testdata/fixtures/pr.json")
	var pr PRData
	json.Unmarshal(data, &pr)

	out := RenderPR(pr, AllSections, nil)
	if !strings.Contains(out, "# PR #42: Fix bug") {
		t.Errorf("missing title in:\n%s", out)
	}
	if !strings.Contains(out, "## Info") {
		t.Errorf("missing Info section in:\n%s", out)
	}
	if !strings.Contains(out, "## Description") {
		t.Errorf("missing Description section in:\n%s", out)
	}
}
