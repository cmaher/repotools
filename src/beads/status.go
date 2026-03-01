package beads

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"repotools/src/runner"
)

type Dependency struct {
	Type        string `json:"type"`
	DependsOnID string `json:"depends_on_id"`
}

type BeadItem struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	IssueType    string       `json:"issue_type"`
	Status       string       `json:"status"`
	Priority     int          `json:"priority"`
	Labels       []string     `json:"labels"`
	Dependencies []Dependency `json:"dependencies"`
}

func (b BeadItem) IsEpic() bool {
	if b.IssueType == "epic" {
		return true
	}
	for _, l := range b.Labels {
		if l == "epic" {
			return true
		}
	}
	return false
}

func BuildStatusReport(items []BeadItem, today string) string {
	children := make(map[string][]BeadItem)
	for _, it := range items {
		for _, dep := range it.Dependencies {
			if dep.Type == "parent-child" {
				children[dep.DependsOnID] = append(children[dep.DependsOnID], it)
			}
		}
	}

	epics := make(map[string]BeadItem)
	for _, it := range items {
		if it.IsEpic() && it.Status == "open" {
			epics[it.ID] = it
		}
	}

	childEpicIDs := make(map[string]bool)
	for pid := range epics {
		for _, child := range children[pid] {
			if child.IsEpic() {
				childEpicIDs[child.ID] = true
			}
		}
	}
	var topEpics []BeadItem
	for id, ep := range epics {
		if !childEpicIDs[id] {
			topEpics = append(topEpics, ep)
		}
	}

	type counts struct {
		open, closed int
		subEpics     []BeadItem
	}

	descendantCounts := func(eid string) counts {
		var c counts
		type stackItem struct {
			id     string
			isRoot bool
		}
		stack := []stackItem{{eid, true}}
		for len(stack) > 0 {
			item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, child := range children[item.id] {
				if child.IsEpic() {
					if child.Status == "closed" {
						continue
					}
					if item.isRoot {
						c.subEpics = append(c.subEpics, child)
					}
					stack = append(stack, stackItem{child.ID, false})
				} else {
					if child.Status == "open" {
						c.open++
					} else if child.Status == "closed" {
						c.closed++
					}
				}
			}
		}
		return c
	}

	var out []string
	out = append(out, "# Project Status", "", fmt.Sprintf("Generated: %s", today), "", "## Open Epics")

	pmap := []struct {
		val   int
		label string
	}{{1, "P1"}, {2, "P2"}, {3, "P3"}}

	for _, p := range pmap {
		var pEpics []BeadItem
		for _, e := range topEpics {
			if e.Priority == p.val {
				pEpics = append(pEpics, e)
			}
		}
		if len(pEpics) == 0 {
			continue
		}
		sort.Slice(pEpics, func(i, j int) bool { return pEpics[i].Title < pEpics[j].Title })

		out = append(out, "", fmt.Sprintf("### %s", p.label), "", "| Epic | ID | Open | Closed |", "|------|----|------|--------|")

		for _, ep := range pEpics {
			c := descendantCounts(ep.ID)
			if c.open+c.closed > 0 {
				out = append(out, fmt.Sprintf("| %s | %s | %d | %d |", ep.Title, ep.ID, c.open, c.closed))
			} else {
				out = append(out, fmt.Sprintf("| %s | %s | — | — |", ep.Title, ep.ID))
			}
			sort.Slice(c.subEpics, func(i, j int) bool { return c.subEpics[i].Title < c.subEpics[j].Title })
			for _, sub := range c.subEpics {
				sc := descendantCounts(sub.ID)
				if sc.open+sc.closed > 0 {
					out = append(out, fmt.Sprintf("| ↳ %s | %s | %d | %d |", sub.Title, sub.ID, sc.open, sc.closed))
				} else {
					out = append(out, fmt.Sprintf("| ↳ %s | %s | — | — |", sub.Title, sub.ID))
				}
			}
		}
	}

	return strings.Join(out, "\n")
}

func findBeadsJSONL() (string, error) {
	dir := ".beads"
	jsonl := filepath.Join(dir, "issues.jsonl")
	if _, err := os.Stat(jsonl); err != nil {
		return "", fmt.Errorf("no .beads/issues.jsonl found: %w", err)
	}
	return jsonl, nil
}

func parseJSONL(path string) ([]BeadItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var items []BeadItem
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var item BeadItem
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, scanner.Err()
}

func RunBeadStatus(w io.Writer) error {
	// Sync DB → JSONL first
	runner.RunNoCheck([]string{"br", "sync", "--flush-only"})

	path, err := findBeadsJSONL()
	if err != nil {
		return err
	}

	items, err := parseJSONL(path)
	if err != nil {
		return fmt.Errorf("parsing beads JSONL: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	fmt.Fprint(w, BuildStatusReport(items, today))
	return nil
}
