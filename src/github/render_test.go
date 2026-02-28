package github

import (
	"strings"
	"testing"
)

func TestFmtTime_Valid(t *testing.T) {
	got := FmtTime("2024-01-15T10:30:00Z")
	want := "2024-01-15 10:30 UTC"
	if got != want {
		t.Errorf("FmtTime = %q, want %q", got, want)
	}
}

func TestFmtTime_Empty(t *testing.T) {
	if got := FmtTime(""); got != "" {
		t.Errorf("FmtTime empty = %q, want empty", got)
	}
}

func TestRenderInfo(t *testing.T) {
	data := PRData{
		Number:       42,
		Title:        "Fix bug",
		State:        "OPEN",
		Author:       Author{Login: "alice"},
		BaseRefName:  "main",
		HeadRefName:  "fix-bug",
		Mergeable:    "MERGEABLE",
		Additions:    10,
		Deletions:    3,
		ChangedFiles: 2,
		URL:          "https://github.com/org/repo/pull/42",
		CreatedAt:    "2024-01-15T10:30:00Z",
		UpdatedAt:    "2024-01-16T12:00:00Z",
	}
	out := RenderInfo(data)
	if !strings.Contains(out, "**#42** Fix bug") {
		t.Errorf("missing title line in:\n%s", out)
	}
	if !strings.Contains(out, "Author: alice") {
		t.Errorf("missing author in:\n%s", out)
	}
	if !strings.Contains(out, "+10 -3 across 2 files") {
		t.Errorf("missing stats in:\n%s", out)
	}
}

func TestRenderInfo_Draft(t *testing.T) {
	data := PRData{Number: 1, Title: "WIP", IsDraft: true, State: "OPEN", Author: Author{Login: "bob"}}
	out := RenderInfo(data)
	if !strings.Contains(out, "(DRAFT)") {
		t.Errorf("missing DRAFT marker in:\n%s", out)
	}
}

func TestRenderBody_Empty(t *testing.T) {
	data := PRData{Body: ""}
	if got := RenderBody(data); got != "(empty)" {
		t.Errorf("RenderBody empty = %q, want '(empty)'", got)
	}
}

func TestRenderBody_Content(t *testing.T) {
	data := PRData{Body: "  some description  "}
	if got := RenderBody(data); got != "some description" {
		t.Errorf("RenderBody = %q, want trimmed", got)
	}
}

func TestRenderComments_Empty(t *testing.T) {
	data := PRData{}
	if got := RenderComments(data); got != "(no comments)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderComments(t *testing.T) {
	data := PRData{
		Comments: []Comment{
			{Author: Author{Login: "alice"}, CreatedAt: "2024-01-15T10:30:00Z", Body: "looks good"},
			{Author: Author{Login: "bob"}, CreatedAt: "2024-01-16T11:00:00Z", Body: "agreed"},
		},
	}
	out := RenderComments(data)
	if !strings.Contains(out, "**alice**") {
		t.Errorf("missing alice in:\n%s", out)
	}
	if !strings.Contains(out, "---") {
		t.Errorf("missing separator in:\n%s", out)
	}
}

func TestRenderReviews_Empty(t *testing.T) {
	if got := RenderReviews(PRData{}); got != "(no reviews)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderReviews(t *testing.T) {
	data := PRData{
		Reviews: []Review{
			{Author: Author{Login: "alice"}, State: "APPROVED", SubmittedAt: "2024-01-15T10:30:00Z", Body: "lgtm"},
		},
	}
	out := RenderReviews(data)
	if !strings.Contains(out, "APPROVED") {
		t.Errorf("missing state in:\n%s", out)
	}
	if !strings.Contains(out, "lgtm") {
		t.Errorf("missing body in:\n%s", out)
	}
}

func TestRenderReviewComments_Empty(t *testing.T) {
	if got := RenderReviewComments(nil); got != "(no inline comments)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderReviewComments(t *testing.T) {
	line := 42
	comments := []ReviewComment{
		{
			User:      Author{Login: "alice"},
			Path:      "main.go",
			Line:      &line,
			CreatedAt: "2024-01-15T10:30:00Z",
			Body:      "nit: rename this",
			DiffHunk:  "+line1\n+line2\n+line3\n+line4",
		},
	}
	out := RenderReviewComments(comments)
	if !strings.Contains(out, "`main.go:42`") {
		t.Errorf("missing file:line in:\n%s", out)
	}
	if !strings.Contains(out, "```diff") {
		t.Errorf("missing diff hunk in:\n%s", out)
	}
}

func TestRenderReviewComments_Reply(t *testing.T) {
	replyTo := 99
	comments := []ReviewComment{
		{
			User:        Author{Login: "bob"},
			Path:        "main.go",
			CreatedAt:   "2024-01-15T10:30:00Z",
			Body:        "done",
			InReplyToID: &replyTo,
		},
	}
	out := RenderReviewComments(comments)
	if !strings.Contains(out, "(reply)") {
		t.Errorf("missing reply note in:\n%s", out)
	}
	if strings.Contains(out, "```diff") {
		t.Errorf("reply should not have diff hunk:\n%s", out)
	}
}

func TestRenderChecks_Empty(t *testing.T) {
	if got := RenderChecks(PRData{}); got != "(no checks)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderChecks(t *testing.T) {
	data := PRData{
		StatusCheckRollup: []CheckStatus{
			{Name: "ci/build", Conclusion: "SUCCESS"},
			{Name: "ci/lint", Conclusion: "FAILURE"},
		},
	}
	out := RenderChecks(data)
	if !strings.Contains(out, "SUCCESS") || !strings.Contains(out, "ci/build") {
		t.Errorf("missing check in:\n%s", out)
	}
}

func TestRenderFiles_Empty(t *testing.T) {
	if got := RenderFiles(PRData{}); got != "(no files)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderFiles(t *testing.T) {
	data := PRData{
		Files: []FileChange{{Path: "main.go", Additions: 10, Deletions: 2}},
	}
	out := RenderFiles(data)
	if !strings.Contains(out, "+10") || !strings.Contains(out, "main.go") {
		t.Errorf("missing file info in:\n%s", out)
	}
}

func TestRenderCommits_Empty(t *testing.T) {
	if got := RenderCommits(PRData{}); got != "(no commits)" {
		t.Errorf("got %q", got)
	}
}

func TestRenderCommits(t *testing.T) {
	data := PRData{
		Commits: []Commit{
			{Oid: "abc1234567890", MessageHeadline: "Fix bug", Authors: []Author{{Login: "alice"}}},
		},
	}
	out := RenderCommits(data)
	if !strings.Contains(out, "abc1234") {
		t.Errorf("missing truncated oid in:\n%s", out)
	}
	if !strings.Contains(out, "Fix bug") {
		t.Errorf("missing message in:\n%s", out)
	}
}
