package beads

import (
	"fmt"
	"io"

	"repotools/src/runner"
)

func FormatMultiBeadHeader(w io.Writer, id string) {
	fmt.Fprintf(w, "=== %s ===\n", id)
}

func RunMultiBead(w io.Writer, ids []string) error {
	for _, id := range ids {
		FormatMultiBeadHeader(w, id)
		r, _ := runner.RunNoCheck([]string{"br", "show", id})
		if r != nil {
			fmt.Fprint(w, r.Stdout)
			if r.Stderr != "" {
				fmt.Fprint(w, r.Stderr)
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}
