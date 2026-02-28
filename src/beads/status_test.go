package beads

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func loadTestItems(t *testing.T) []BeadItem {
	t.Helper()
	data, err := os.ReadFile("../../testdata/beads/sample_list.json")
	if err != nil {
		t.Fatal(err)
	}
	var items []BeadItem
	if err := json.Unmarshal(data, &items); err != nil {
		t.Fatal(err)
	}
	return items
}

func TestBuildStatusReport(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "# Project Status") {
		t.Errorf("missing header in:\n%s", report)
	}
	if !strings.Contains(report, "Epic One") {
		t.Errorf("missing Epic One in:\n%s", report)
	}
	if !strings.Contains(report, "↳ Sub Epic") {
		t.Errorf("missing sub-epic in:\n%s", report)
	}
}

func TestBuildStatusReport_Counts(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "| Epic One | E1 | 2 | 1 |") {
		t.Errorf("wrong counts for Epic One in:\n%s", report)
	}
	if !strings.Contains(report, "| ↳ Sub Epic | E2 | 1 | 0 |") {
		t.Errorf("wrong counts for Sub Epic in:\n%s", report)
	}
}

func TestRunBeadStatus_WritesToWriter(t *testing.T) {
	items := loadTestItems(t)
	var buf bytes.Buffer
	report := BuildStatusReport(items, "2026-02-28")
	buf.WriteString(report)

	if !strings.Contains(buf.String(), "Generated: 2026-02-28") {
		t.Errorf("missing date in output")
	}
}

func TestBuildStatusReport_LabelEpic(t *testing.T) {
	items := []BeadItem{
		{ID: "X1", Title: "Labeled Epic", IssueType: "task", Status: "open", Priority: 1, Labels: []string{"epic"}},
		{ID: "X2", Title: "Child Task", IssueType: "task", Status: "open", Priority: 2, Dependencies: []Dependency{{Type: "parent-child", DependsOnID: "X1"}}},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if !strings.Contains(report, "Labeled Epic") {
		t.Errorf("label-based epic missing in:\n%s", report)
	}
	if !strings.Contains(report, "| Labeled Epic | X1 | 1 | 0 |") {
		t.Errorf("wrong counts for labeled epic in:\n%s", report)
	}
}

func TestBuildStatusReport_ClosedEpicExcluded(t *testing.T) {
	items := []BeadItem{
		{ID: "C1", Title: "Closed Epic", IssueType: "epic", Status: "closed", Priority: 1},
		{ID: "C2", Title: "Child", IssueType: "task", Status: "open", Priority: 2, Dependencies: []Dependency{{Type: "parent-child", DependsOnID: "C1"}}},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if strings.Contains(report, "Closed Epic") {
		t.Errorf("closed epic should not appear in:\n%s", report)
	}
}

func TestBuildStatusReport_ClosedSubEpicExcluded(t *testing.T) {
	items := []BeadItem{
		{ID: "P1", Title: "Parent Epic", IssueType: "epic", Status: "open", Priority: 1},
		{ID: "S1", Title: "Closed Sub", IssueType: "epic", Status: "closed", Priority: 1, Dependencies: []Dependency{{Type: "parent-child", DependsOnID: "P1"}}},
		{ID: "S2", Title: "Open Sub", IssueType: "epic", Status: "open", Priority: 1, Dependencies: []Dependency{{Type: "parent-child", DependsOnID: "P1"}}},
		{ID: "T1", Title: "Task", IssueType: "task", Status: "open", Priority: 2, Dependencies: []Dependency{{Type: "parent-child", DependsOnID: "S2"}}},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if strings.Contains(report, "Closed Sub") {
		t.Errorf("closed sub-epic should not appear in:\n%s", report)
	}
	if !strings.Contains(report, "Open Sub") {
		t.Errorf("open sub-epic missing in:\n%s", report)
	}
}
