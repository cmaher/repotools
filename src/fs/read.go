package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadLines(w io.Writer, path string, start, end int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s: Is a directory", path)
	}

	scanner := bufio.NewScanner(f)
	lineno := 0
	for scanner.Scan() {
		lineno++
		if start > 0 && lineno < start {
			continue
		}
		if end > 0 && lineno > end {
			break
		}
		fmt.Fprintf(w, "%6d\t%s\n", lineno, scanner.Text())
	}
	return scanner.Err()
}
