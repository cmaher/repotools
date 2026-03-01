package tickets

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Ticket struct {
	ID       string
	Title    string
	Type     string
	Status   string
	Priority int
	Parent   string
	Tags     []string
}

func (t Ticket) IsEpic() bool {
	return t.Type == "epic"
}

// parseTicketFile reads a single .md ticket file and extracts frontmatter fields.
func parseTicketFile(path string) (Ticket, error) {
	f, err := os.Open(path)
	if err != nil {
		return Ticket{}, err
	}
	defer f.Close()

	var tk Ticket
	scanner := bufio.NewScanner(f)

	// Expect opening ---
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return Ticket{}, fmt.Errorf("%s: missing frontmatter", path)
	}

	// Read frontmatter lines until closing ---
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		key, val, ok := parseYAMLLine(line)
		if !ok {
			continue
		}
		switch key {
		case "id":
			tk.ID = val
		case "status":
			tk.Status = val
		case "type":
			tk.Type = val
		case "priority":
			if n, err := strconv.Atoi(val); err == nil {
				tk.Priority = n
			}
		case "parent":
			tk.Parent = val
		case "tags":
			tk.Tags = parseYAMLList(val)
		}
	}

	// First markdown heading is the title
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			tk.Title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	return tk, scanner.Err()
}

// parseYAMLLine does minimal key: value extraction from frontmatter.
func parseYAMLLine(line string) (key, val string, ok bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	val = strings.TrimSpace(line[idx+1:])
	return key, val, true
}

// parseYAMLList handles both [a, b] and bare comma-separated values.
func parseYAMLList(val string) []string {
	val = strings.TrimPrefix(val, "[")
	val = strings.TrimSuffix(val, "]")
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// LoadTickets reads all .md files from a directory.
func LoadTickets(dir string) ([]Ticket, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading tickets dir %s: %w", dir, err)
	}

	var tickets []Ticket
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		tk, err := parseTicketFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		tickets = append(tickets, tk)
	}
	return tickets, nil
}

func BuildStatusReport(items []Ticket, today string) string {
	children := make(map[string][]Ticket)
	for _, it := range items {
		if it.Parent != "" {
			children[it.Parent] = append(children[it.Parent], it)
		}
	}

	epics := make(map[string]Ticket)
	for _, it := range items {
		if it.IsEpic() && it.Status != "closed" {
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
	var topEpics []Ticket
	for id, ep := range epics {
		if !childEpicIDs[id] {
			topEpics = append(topEpics, ep)
		}
	}

	// Collect direct sub-epics for each epic
	subEpicsOf := func(eid string) []Ticket {
		var subs []Ticket
		for _, child := range children[eid] {
			if child.IsEpic() && child.Status != "closed" {
				subs = append(subs, child)
			}
		}
		sort.Slice(subs, func(i, j int) bool { return subs[i].Title < subs[j].Title })
		return subs
	}

	// Count open/closed non-epic descendants
	type counts struct{ open, closed int }
	descendantCounts := func(eid string) counts {
		var c counts
		stack := []string{eid}
		for len(stack) > 0 {
			id := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, child := range children[id] {
				if child.IsEpic() {
					stack = append(stack, child.ID)
				} else if child.Status == "closed" {
					c.closed++
				} else {
					c.open++
				}
			}
		}
		return c
	}

	var totalOpen, totalClosed int
	for _, it := range items {
		if it.Status == "closed" {
			totalClosed++
		} else {
			totalOpen++
		}
	}

	var out []string
	out = append(out, fmt.Sprintf("Project Status — %s  (%d open, %d closed)", today, totalOpen, totalClosed), "")

	pmap := []struct {
		val   int
		label string
	}{{1, "P1"}, {2, "P2"}, {3, "P3"}}

	for _, p := range pmap {
		var pEpics []Ticket
		for _, e := range topEpics {
			if e.Priority == p.val {
				pEpics = append(pEpics, e)
			}
		}
		if len(pEpics) == 0 {
			continue
		}
		sort.Slice(pEpics, func(i, j int) bool { return pEpics[i].Title < pEpics[j].Title })

		out = append(out, fmt.Sprintf("[%s]", p.label))

		for _, ep := range pEpics {
			c := descendantCounts(ep.ID)
			out = append(out, fmt.Sprintf("  %-12s %-7s %s", ep.ID, fmt.Sprintf("%d/%d", c.open, c.closed), ep.Title))

			for _, sub := range subEpicsOf(ep.ID) {
				sc := descendantCounts(sub.ID)
				out = append(out, fmt.Sprintf("    %-12s %-7s %s", sub.ID, fmt.Sprintf("%d/%d", sc.open, sc.closed), sub.Title))
			}
		}
		out = append(out, "")
	}

	return strings.Join(out, "\n")
}

func findTicketsDir() (string, error) {
	dir := ".tickets"
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("no .tickets/ directory found")
	}
	return dir, nil
}

func RunTicketStatus(w io.Writer) error {
	dir, err := findTicketsDir()
	if err != nil {
		return err
	}

	items, err := LoadTickets(dir)
	if err != nil {
		return fmt.Errorf("loading tickets: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	fmt.Fprint(w, BuildStatusReport(items, today))
	return nil
}
