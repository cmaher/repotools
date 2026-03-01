package cli

import (
	"fmt"
	"os"
	"strconv"

	"repotools/src/fs"

	"github.com/spf13/cobra"
)

func newMultiFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "multi-find <head_count> [find_opts...] path1 [path2 ...]",
		Aliases: []string{"mf"},
		Short:   "Find files in multiple directories with truncated output",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			headCount, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("head_count must be an integer: %w", err)
			}
			rest := args[1:]

			var findOpts, paths []string
			for _, a := range rest {
				if _, err := os.Stat(a); err == nil {
					paths = append(paths, a)
				} else {
					findOpts = append(findOpts, a)
				}
			}
			if len(paths) == 0 {
				return fmt.Errorf("no valid paths found")
			}

			fs.MultiFind(os.Stdout, headCount, findOpts, paths)
			return nil
		},
	}
}
