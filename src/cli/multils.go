package cli

import (
	"os"

	"repotools/src/fs"

	"github.com/spf13/cobra"
)

func newMultiLSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "multi-ls dir1 dir2 ...",
		Short: "List contents of multiple directories",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fs.MultiLS(os.Stdout, args)
		},
	}
}
