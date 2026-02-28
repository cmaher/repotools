package runner

import (
	"strings"
	"testing"
)

func TestRun_CapturesOutput(t *testing.T) {
	result, err := Run([]string{"echo", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(result.Stdout); got != "hello" {
		t.Errorf("stdout = %q, want %q", got, "hello")
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	_, err := Run([]string{"false"})
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
}

func TestRunNoCheck_NonZeroExit(t *testing.T) {
	result, err := RunNoCheck([]string{"false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code")
	}
}
