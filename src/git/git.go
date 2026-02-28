package git

import (
	"fmt"
	"io"
	"strings"

	"repotools/src/runner"
)

func MergeBase(base string) (string, error) {
	r, err := runner.Run([]string{"git", "merge-base", "HEAD", base})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(r.Stdout), nil
}

func Status(w io.Writer) error {
	branch, err := runner.Run([]string{"git", "branch", "--show-current"})
	if err != nil {
		return err
	}
	status, _ := runner.RunNoCheck([]string{"git", "status", "--short"})

	fmt.Fprintf(w, "Branch: %s\n", strings.TrimSpace(branch.Stdout))
	fmt.Fprintln(w, "---")
	fmt.Fprint(w, status.Stdout)
	return nil
}
