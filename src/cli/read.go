package cli

import (
	"os"
	"strconv"

	"repotools/src/fs"

	"github.com/spf13/cobra"
)

func newReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "read <file> [start] [end]",
		Short: "Print numbered lines from a file (optional range)",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			start, end := 0, 0
			if len(args) >= 2 {
				var err error
				start, err = strconv.Atoi(args[1])
				if err != nil {
					return err
				}
			}
			if len(args) >= 3 {
				var err error
				end, err = strconv.Atoi(args[2])
				if err != nil {
					return err
				}
			}
			return fs.ReadLines(os.Stdout, path, start, end)
		},
	}
}
