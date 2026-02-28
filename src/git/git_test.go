package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "checkout", "-b", "main"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %s: %v", args, out, err)
		}
	}

	f := dir + "/file.txt"
	os.WriteFile(f, []byte("hello\n"), 0644)
	for _, args := range [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %s: %v", args, out, err)
		}
	}

	return dir
}

func TestMergeBase(t *testing.T) {
	dir := setupGitRepo(t)
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	mb, err := MergeBase("main")
	if err != nil {
		t.Fatalf("MergeBase: %v", err)
	}
	if len(mb) < 7 {
		t.Errorf("MergeBase returned %q, want a commit hash", mb)
	}
}

func TestStatus(t *testing.T) {
	dir := setupGitRepo(t)
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	var buf bytes.Buffer
	err := Status(&buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Branch: main") {
		t.Errorf("output missing branch, got: %s", out)
	}
	if !strings.Contains(out, "---") {
		t.Errorf("output missing separator, got: %s", out)
	}
}
