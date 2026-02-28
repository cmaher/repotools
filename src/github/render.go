package github

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func FmtTime(ts string) string {
	if ts == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.UTC().Format("2006-01-02 15:04 UTC")
}

func RenderInfo(data PRData) string {
	author := data.Author.Login
	if author == "" {
		author = "unknown"
	}

	labels := "none"
	if len(data.Labels) > 0 {
		names := make([]string, len(data.Labels))
		for i, l := range data.Labels {
			names[i] = l.Name
		}
		labels = strings.Join(names, ", ")
	}

	assignees := "none"
	if len(data.Assignees) > 0 {
		logins := make([]string, len(data.Assignees))
		for i, a := range data.Assignees {
			logins[i] = a.Login
		}
		assignees = strings.Join(logins, ", ")
	}

	reviewRequests := "none"
	if len(data.ReviewRequests) > 0 {
		names := make([]string, len(data.ReviewRequests))
		for i, r := range data.ReviewRequests {
			if r.Login != "" {
				names[i] = r.Login
			} else {
				names[i] = r.Name
			}
		}
		reviewRequests = strings.Join(names, ", ")
	}

	draft := ""
	if data.IsDraft {
		draft = " (DRAFT)"
	}

	reviewDecision := data.ReviewDecision
	if reviewDecision == "" {
		reviewDecision = "PENDING"
	}

	mergeable := data.Mergeable
	if mergeable == "" {
		mergeable = "?"
	}

	lines := []string{
		fmt.Sprintf("**#%d** %s%s", data.Number, data.Title, draft),
		fmt.Sprintf("State: %s | Review: %s | Mergeable: %s", data.State, reviewDecision, mergeable),
		fmt.Sprintf("Author: %s | Assignees: %s", author, assignees),
		fmt.Sprintf("Base: %s <- Head: %s", data.BaseRefName, data.HeadRefName),
		fmt.Sprintf("Labels: %s | Reviewers requested: %s", labels, reviewRequests),
		fmt.Sprintf("+%d -%d across %d files", data.Additions, data.Deletions, data.ChangedFiles),
		fmt.Sprintf("Created: %s | Updated: %s", FmtTime(data.CreatedAt), FmtTime(data.UpdatedAt)),
	}

	if data.MergedAt != "" {
		lines = append(lines, fmt.Sprintf("Merged: %s", FmtTime(data.MergedAt)))
	}
	if data.ClosedAt != "" && data.State != "MERGED" {
		lines = append(lines, fmt.Sprintf("Closed: %s", FmtTime(data.ClosedAt)))
	}
	lines = append(lines, fmt.Sprintf("URL: %s", data.URL))

	return strings.Join(lines, "\n")
}

func RenderBody(data PRData) string {
	body := strings.TrimSpace(data.Body)
	if body == "" {
		return "(empty)"
	}
	return body
}

func RenderComments(data PRData) string {
	if len(data.Comments) == 0 {
		return "(no comments)"
	}
	parts := make([]string, len(data.Comments))
	for i, c := range data.Comments {
		author := c.Author.Login
		if author == "" {
			author = "unknown"
		}
		parts[i] = fmt.Sprintf("**%s** (%s):\n%s", author, FmtTime(c.CreatedAt), strings.TrimSpace(c.Body))
	}
	return strings.Join(parts, "\n\n---\n\n")
}

func RenderReviews(data PRData) string {
	if len(data.Reviews) == 0 {
		return "(no reviews)"
	}
	parts := make([]string, len(data.Reviews))
	for i, r := range data.Reviews {
		author := r.Author.Login
		if author == "" {
			author = "unknown"
		}
		ts := r.SubmittedAt
		if ts == "" {
			ts = r.CreatedAt
		}
		state := r.State
		if state == "" {
			state = "?"
		}
		entry := fmt.Sprintf("**%s** — %s (%s)", author, state, FmtTime(ts))
		body := strings.TrimSpace(r.Body)
		if body != "" {
			entry += "\n" + body
		}
		parts[i] = entry
	}
	return strings.Join(parts, "\n\n---\n\n")
}

func RenderReviewComments(comments []ReviewComment) string {
	if len(comments) == 0 {
		return "(no inline comments)"
	}
	parts := make([]string, len(comments))
	for i, c := range comments {
		author := c.User.Login
		if author == "" {
			author = "unknown"
		}
		path := c.Path
		if path == "" {
			path = "?"
		}
		line := "?"
		if c.OriginalLine != nil {
			line = fmt.Sprintf("%d", *c.OriginalLine)
		} else if c.Line != nil {
			line = fmt.Sprintf("%d", *c.Line)
		}

		replyNote := ""
		if c.InReplyToID != nil {
			replyNote = " (reply)"
		}

		entry := fmt.Sprintf("**%s**%s on `%s:%s` (%s):", author, replyNote, path, line, FmtTime(c.CreatedAt))

		if c.DiffHunk != "" && c.InReplyToID == nil {
			hunkLines := strings.Split(c.DiffHunk, "\n")
			start := len(hunkLines) - 3
			if start < 0 {
				start = 0
			}
			entry += "\n```diff\n" + strings.Join(hunkLines[start:], "\n") + "\n```"
		}

		body := strings.TrimSpace(c.Body)
		if body != "" {
			entry += "\n" + body
		}
		parts[i] = entry
	}
	return strings.Join(parts, "\n\n---\n\n")
}

func RenderChecks(data PRData) string {
	if len(data.StatusCheckRollup) == 0 {
		return "(no checks)"
	}
	lines := make([]string, len(data.StatusCheckRollup))
	for i, c := range data.StatusCheckRollup {
		name := c.Name
		if name == "" {
			name = c.Context
		}
		if name == "" {
			name = "?"
		}
		status := c.Conclusion
		if status == "" {
			status = c.Status
		}
		if status == "" {
			status = c.State
		}
		if status == "" {
			status = "?"
		}
		lines[i] = fmt.Sprintf("  %-20s %s", status, name)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func RenderFiles(data PRData) string {
	if len(data.Files) == 0 {
		return "(no files)"
	}
	lines := make([]string, len(data.Files))
	for i, f := range data.Files {
		lines[i] = fmt.Sprintf("  +%-4d -%-4d %s", f.Additions, f.Deletions, f.Path)
	}
	return strings.Join(lines, "\n")
}

func RenderCommits(data PRData) string {
	if len(data.Commits) == 0 {
		return "(no commits)"
	}
	lines := make([]string, len(data.Commits))
	for i, c := range data.Commits {
		oid := c.Oid
		if len(oid) > 7 {
			oid = oid[:7]
		}
		if oid == "" {
			oid = "???????"
		}
		authors := make([]string, len(c.Authors))
		for j, a := range c.Authors {
			if a.Login != "" {
				authors[j] = a.Login
			} else if a.Name != "" {
				authors[j] = a.Name
			} else {
				authors[j] = "?"
			}
		}
		lines[i] = fmt.Sprintf("  %s %s (%s)", oid, c.MessageHeadline, strings.Join(authors, ", "))
	}
	return strings.Join(lines, "\n")
}
