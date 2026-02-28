package fs

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func MultiFind(w io.Writer, headCount int, findOpts []string, paths []string) {
	for _, p := range paths {
		fmt.Fprintf(w, "==> %s <==\n", p)

		args := append([]string{p}, findOpts...)
		cmd := exec.Command("find", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(w, err)
			fmt.Fprintln(w, "---")
			continue
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintln(w, err)
			fmt.Fprintln(w, "---")
			continue
		}

		scanner := bufio.NewScanner(stdout)
		count := 0
		for scanner.Scan() {
			if count >= headCount {
				break
			}
			fmt.Fprintln(w, scanner.Text())
			count++
		}

		for scanner.Scan() {
		}

		cmd.Wait()
		fmt.Fprintln(w, "---")
	}
}
