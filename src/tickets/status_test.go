package tickets

import (
	"strings"
	"testing"
)

func loadTestItems(t *testing.T) []Ticket {
	t.Helper()
	items, err := LoadTickets("../../testdata/tickets")
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func TestLoadTickets(t *testing.T) {
	items := loadTestItems(t)
	if len(items) != 7 {
		t.Fatalf("expected 7 tickets, got %d", len(items))
	}
}

func TestBuildStatusReport(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "Project Status") {
		t.Errorf("missing header in:\n%s", report)
	}
	if !strings.Contains(report, "Epic One") {
		t.Errorf("missing Epic One in:\n%s", report)
	}
	if !strings.Contains(report, "Sub Epic") {
		t.Errorf("missing sub-epic in:\n%s", report)
	}
}

func TestBuildStatusReport_EpicRows(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "E1") || !strings.Contains(report, "Epic One") {
		t.Errorf("missing Epic One row in:\n%s", report)
	}
	if !strings.Contains(report, "E2") || !strings.Contains(report, "Sub Epic") {
		t.Errorf("missing Sub Epic row in:\n%s", report)
	}
}

func TestBuildStatusReport_Date(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "2026-02-28") {
		t.Errorf("missing date in output")
	}
}

func TestBuildStatusReport_ClosedEpicExcluded(t *testing.T) {
	items := []Ticket{
		{ID: "C1", Title: "Closed Epic", Type: "epic", Status: "closed", Priority: 1},
		{ID: "C2", Title: "Child", Type: "task", Status: "open", Priority: 2, Parent: "C1"},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if strings.Contains(report, "Closed Epic") {
		t.Errorf("closed epic should not appear in:\n%s", report)
	}
}

func TestBuildStatusReport_ClosedSubEpicExcluded(t *testing.T) {
	items := []Ticket{
		{ID: "P1", Title: "Parent Epic", Type: "epic", Status: "open", Priority: 1},
		{ID: "S1", Title: "Closed Sub", Type: "epic", Status: "closed", Priority: 1, Parent: "P1"},
		{ID: "S2", Title: "Open Sub", Type: "epic", Status: "open", Priority: 1, Parent: "P1"},
		{ID: "T1", Title: "Task", Type: "task", Status: "open", Priority: 2, Parent: "S2"},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if strings.Contains(report, "Closed Sub") {
		t.Errorf("closed sub-epic should not appear in:\n%s", report)
	}
	if !strings.Contains(report, "Open Sub") {
		t.Errorf("open sub-epic missing in:\n%s", report)
	}
}

func TestBuildStatusReport_OrphanedTickets(t *testing.T) {
	items := loadTestItems(t)
	report := BuildStatusReport(items, "2026-02-28")

	if !strings.Contains(report, "[Orphaned]") {
		t.Errorf("missing Orphaned section in:\n%s", report)
	}
	if !strings.Contains(report, "Orphan Feature") {
		t.Errorf("missing Orphan Feature in:\n%s", report)
	}
	if !strings.Contains(report, "Orphan Task") {
		t.Errorf("missing Orphan Task in:\n%s", report)
	}
}

func TestBuildStatusReport_OrphanedClosedExcluded(t *testing.T) {
	items := []Ticket{
		{ID: "O1", Title: "Open Orphan", Type: "task", Status: "open", Priority: 2},
		{ID: "O2", Title: "Closed Orphan", Type: "task", Status: "closed", Priority: 2},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if !strings.Contains(report, "Open Orphan") {
		t.Errorf("missing open orphan in:\n%s", report)
	}
	if strings.Contains(report, "Closed Orphan") {
		t.Errorf("closed orphan should not appear in:\n%s", report)
	}
}

func TestBuildStatusReport_NoOrphansWhenAllParented(t *testing.T) {
	items := []Ticket{
		{ID: "E1", Title: "Epic", Type: "epic", Status: "open", Priority: 1},
		{ID: "T1", Title: "Task", Type: "task", Status: "open", Priority: 2, Parent: "E1"},
	}
	report := BuildStatusReport(items, "2026-02-28")
	if strings.Contains(report, "[Orphaned]") {
		t.Errorf("should not have Orphaned section when all tickets are parented:\n%s", report)
	}
}

func TestParseTicketFile(t *testing.T) {
	tk, err := parseTicketFile("../../testdata/tickets/T1.md")
	if err != nil {
		t.Fatal(err)
	}
	if tk.ID != "T1" {
		t.Errorf("expected ID T1, got %s", tk.ID)
	}
	if tk.Title != "Task One" {
		t.Errorf("expected title 'Task One', got %q", tk.Title)
	}
	if tk.Parent != "E1" {
		t.Errorf("expected parent E1, got %s", tk.Parent)
	}
	if tk.Type != "task" {
		t.Errorf("expected type task, got %s", tk.Type)
	}
	if tk.Status != "open" {
		t.Errorf("expected status open, got %s", tk.Status)
	}
	if tk.Priority != 2 {
		t.Errorf("expected priority 2, got %d", tk.Priority)
	}
}
