package metrics

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var DefaultTestMarkers = map[string]string{
	".rs": `^#\[cfg\(test\)\]`,
}

func CountLOC(path string, marker string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var markerRe *regexp.Regexp
	if marker != "" {
		markerRe, err = regexp.Compile(marker)
		if err != nil {
			return 0, err
		}
	}

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if markerRe != nil && markerRe.MatchString(scanner.Text()) {
			break
		}
		count++
	}
	return count, scanner.Err()
}

func CountLOCAutoDetect(path string, explicitMarker string) (int, error) {
	marker := explicitMarker
	if marker == "" {
		ext := filepath.Ext(path)
		marker = DefaultTestMarkers[ext]
	}
	return CountLOC(path, marker)
}

func ResolveFiles(paths []string, globPattern string, excludePattern string) ([]string, error) {
	var excludeRe *regexp.Regexp
	if excludePattern != "" {
		var err error
		excludeRe, err = regexp.Compile(excludePattern)
		if err != nil {
			return nil, err
		}
	}

	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			matches, _ := filepath.Glob(p)
			sort.Strings(matches)
			for _, m := range matches {
				fi, err := os.Stat(m)
				if err == nil && !fi.IsDir() {
					files = append(files, m)
				}
			}
			continue
		}

		if !info.IsDir() {
			if globPattern == "" || matchGlob(p, globPattern) {
				files = append(files, p)
			}
			continue
		}

		pattern := globPattern
		if pattern == "" {
			pattern = "*"
		}
		err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			if matchGlob(path, pattern) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if excludeRe != nil {
		var filtered []string
		for _, f := range files {
			if !excludeRe.MatchString(f) {
				filtered = append(filtered, f)
			}
		}
		files = filtered
	}

	return files, nil
}

func matchGlob(path, pattern string) bool {
	matched, _ := filepath.Match(pattern, filepath.Base(path))
	return matched
}

func RunLOC(w io.Writer, paths []string, globPattern, excludePattern, marker string) error {
	files, err := ResolveFiles(paths, globPattern, excludePattern)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("No files found.")
	}

	total := 0
	for _, f := range files {
		n, err := CountLOCAutoDetect(f, marker)
		if err != nil {
			return err
		}
		total += n
		fmt.Fprintf(w, "%6d %s\n", n, f)
	}

	if len(files) > 1 {
		fmt.Fprintf(w, "%6d total\n", total)
	}
	return nil
}
