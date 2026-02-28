package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"repotools/src/runner"
)

var AllSections = []string{"info", "body", "comments", "reviews", "review-comments", "checks", "files", "commits"}

var GHPRFields = "number,title,body,state,author,baseRefName,headRefName,headRefOid,url," +
	"labels,assignees,reviewRequests,createdAt,updatedAt,mergedAt,closedAt," +
	"additions,deletions,changedFiles,mergeable,reviewDecision,isDraft," +
	"comments,reviews,commits,files,statusCheckRollup"

func ValidateSections(only string) ([]string, error) {
	valid := make(map[string]bool)
	for _, s := range AllSections {
		valid[s] = true
	}
	sections := strings.Split(only, ",")
	for i := range sections {
		sections[i] = strings.TrimSpace(sections[i])
	}
	var bad []string
	for _, s := range sections {
		if !valid[s] {
			bad = append(bad, s)
		}
	}
	if len(bad) > 0 {
		return nil, fmt.Errorf("unknown sections: %s\nAvailable: %s", strings.Join(bad, ", "), strings.Join(AllSections, ", "))
	}
	return sections, nil
}

func FilterSections(base []string, only, exclude string) []string {
	if only != "" {
		sections := strings.Split(only, ",")
		for i := range sections {
			sections[i] = strings.TrimSpace(sections[i])
		}
		return sections
	}
	if exclude != "" {
		excluded := make(map[string]bool)
		for _, s := range strings.Split(exclude, ",") {
			excluded[strings.TrimSpace(s)] = true
		}
		var result []string
		for _, s := range base {
			if !excluded[s] {
				result = append(result, s)
			}
		}
		return result
	}
	return base
}

func FetchPRData(prArg string) (*PRData, error) {
	args := []string{"gh", "pr", "view"}
	if prArg != "" {
		args = append(args, prArg)
	}
	args = append(args, "--json", GHPRFields)

	r, err := runner.RunNoCheck(args)
	if err != nil {
		return nil, err
	}
	if r.ExitCode != 0 {
		msg := strings.TrimSpace(r.Stderr)
		if msg == "" {
			msg = "failed to fetch PR"
		}
		return nil, fmt.Errorf("%s", msg)
	}

	var data PRData
	if err := json.Unmarshal([]byte(r.Stdout), &data); err != nil {
		return nil, fmt.Errorf("parsing PR JSON: %w", err)
	}
	return &data, nil
}

func FetchReviewComments(repo string, prNumber int) ([]ReviewComment, error) {
	r, err := runner.RunNoCheck([]string{
		"gh", "api", fmt.Sprintf("repos/%s/pulls/%d/comments", repo, prNumber),
		"--paginate", "-q", ".[]",
	})
	if err != nil || r.ExitCode != 0 {
		return nil, nil
	}

	var comments []ReviewComment
	for _, line := range strings.Split(strings.TrimSpace(r.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var c ReviewComment
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			continue
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func GetRepoNWO() (string, error) {
	r, err := runner.Run([]string{"gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner"})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(r.Stdout), nil
}

func RenderPR(data PRData, sections []string, reviewComments []ReviewComment) string {
	type sectionDef struct {
		title    string
		renderer func() string
	}

	renderers := map[string]sectionDef{
		"info":            {"Info", func() string { return RenderInfo(data) }},
		"body":            {"Description", func() string { return RenderBody(data) }},
		"comments":        {"Comments", func() string { return RenderComments(data) }},
		"reviews":         {"Reviews", func() string { return RenderReviews(data) }},
		"review-comments": {"Review Comments (inline)", func() string { return RenderReviewComments(reviewComments) }},
		"checks":          {"Checks", func() string { return RenderChecks(data) }},
		"files":           {"Files Changed", func() string { return RenderFiles(data) }},
		"commits":         {"Commits", func() string { return RenderCommits(data) }},
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# PR #%d: %s", data.Number, data.Title)
	for _, s := range sections {
		def := renderers[s]
		fmt.Fprintf(&sb, "\n\n## %s\n\n%s", def.title, def.renderer())
	}
	return sb.String()
}
