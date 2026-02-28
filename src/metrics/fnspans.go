package metrics

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type FnSpan struct {
	Start int
	End   int
	Name  string
}

var DefaultFnPatterns = map[string]string{
	".rs": `^\s*(?:pub(?:\(crate\))?\s+)?(?:async\s+)?fn\s+(\w+)`,
	".py": `^\s*(?:async\s+)?def\s+(\w+)`,
	".js": `^\s*(?:async\s+)?function\s+(\w+)`,
	".ts": `^\s*(?:async\s+)?function\s+(\w+)`,
	".go": `^func\s+(?:\([^)]*\)\s+)?(\w+)`,
}

func fnPatternForFile(path string, explicit string) (*regexp.Regexp, error) {
	pattern := explicit
	if pattern == "" {
		ext := filepath.Ext(path)
		pattern = DefaultFnPatterns[ext]
	}
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

func ExtractFnSpans(path string, pattern string, after string, include string, exclude string) ([]FnSpan, error) {
	fnRe, err := fnPatternForFile(path, pattern)
	if err != nil {
		return nil, err
	}
	if fnRe == nil {
		return nil, nil
	}

	var afterRe, includeRe, excludeRe *regexp.Regexp
	if after != "" {
		afterRe, err = regexp.Compile(after)
		if err != nil {
			return nil, err
		}
	}
	if include != "" {
		includeRe, err = regexp.Compile(include)
		if err != nil {
			return nil, err
		}
	}
	if exclude != "" {
		excludeRe, err = regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	scanning := afterRe == nil
	type raw struct {
		line int
		name string
	}
	var raws []raw

	for i, line := range lines {
		lineno := i + 1
		if !scanning {
			if afterRe.MatchString(line) {
				scanning = true
			}
			continue
		}
		m := fnRe.FindStringSubmatch(line)
		if m != nil && len(m) > 1 {
			raws = append(raws, raw{lineno, m[1]})
		}
	}

	if len(raws) == 0 {
		return nil, nil
	}

	var spans []FnSpan
	for i, r := range raws {
		end := len(lines)
		if i+1 < len(raws) {
			end = raws[i+1].line - 1
		}
		spans = append(spans, FnSpan{Start: r.line, End: end, Name: r.name})
	}

	if includeRe != nil {
		var filtered []FnSpan
		for _, s := range spans {
			if includeRe.MatchString(s.Name) {
				filtered = append(filtered, s)
			}
		}
		spans = filtered
	}
	if excludeRe != nil {
		var filtered []FnSpan
		for _, s := range spans {
			if !excludeRe.MatchString(s.Name) {
				filtered = append(filtered, s)
			}
		}
		spans = filtered
	}

	return spans, nil
}

func RunFnSpans(w io.Writer, paths []string, globPattern, excludePath, pattern, after, include, exclude string) error {
	files, err := ResolveFiles(paths, globPattern, excludePath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("No files found.")
	}

	multi := len(files) > 1

	for _, f := range files {
		fnRe, _ := fnPatternForFile(f, pattern)
		if fnRe == nil {
			if !multi {
				return fmt.Errorf("No function pattern for %s, use --pattern", f)
			}
			fmt.Fprintf(os.Stderr, "No function pattern for %s, use --pattern\n", f)
			continue
		}

		spans, err := ExtractFnSpans(f, pattern, after, include, exclude)
		if err != nil {
			return err
		}
		if len(spans) == 0 {
			if !multi {
				return fmt.Errorf("No functions found.")
			}
			continue
		}

		if multi {
			fmt.Fprintf(w, "==> %s <==\n", f)
		}
		for _, s := range spans {
			size := s.End - s.Start + 1
			fmt.Fprintf(w, "  %d-%d %s (%d lines)\n", s.Start, s.End, s.Name, size)
		}
		if multi {
			fmt.Fprintln(w)
		}
	}
	return nil
}
