package beads

import (
	"encoding/json"
	"fmt"
	"io"
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

func RunBeadStatus(w io.Writer) error {
	r, err := runner.RunNoCheck([]string{"bd", "list", "-n", "0", "--all", "--json"})
	if err != nil {
		return err
	}

	var items []BeadItem
	stdout := strings.TrimSpace(r.Stdout)
	if stdout != "" {
		if err := json.Unmarshal([]byte(stdout), &items); err != nil {
			return fmt.Errorf("parsing beads JSON: %w", err)
		}
	}

	// Fetch items with "epic" label to populate Labels field
	lr, err := runner.RunNoCheck([]string{"bd", "list", "-n", "0", "--all", "--json", "--label-any", "epic"})
	if err == nil && lr.ExitCode == 0 {
		lout := strings.TrimSpace(lr.Stdout)
		if lout != "" {
			var labeled []BeadItem
			if err := json.Unmarshal([]byte(lout), &labeled); err == nil {
				epicLabeled := make(map[string][]string)
				for _, li := range labeled {
					epicLabeled[li.ID] = li.Labels
				}
				for i := range items {
					if labels, ok := epicLabeled[items[i].ID]; ok {
						items[i].Labels = labels
					}
				}
			}
		}
	}

	today := time.Now().Format("2006-01-02")
	fmt.Fprint(w, BuildStatusReport(items, today))
	return nil
}
