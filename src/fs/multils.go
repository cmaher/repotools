package fs

import (
	"fmt"
	"io"
	"os"
)

func MultiLS(w io.Writer, dirs []string) {
	for _, d := range dirs {
		fmt.Fprintf(w, "==> %s <==\n", d)
		entries, err := os.ReadDir(d)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			for _, e := range entries {
				fmt.Fprintln(w, e.Name())
			}
		}
		fmt.Fprintln(w, "---")
	}
}
