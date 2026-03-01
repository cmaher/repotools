package cli

import (
	"os"

	"repotools/src/metrics"

	"github.com/spf13/cobra"
)

func newLocCmd() *cobra.Command {
	var glob, exclude, marker string

	cmd := &cobra.Command{
		Use:     "loc paths...",
		Aliases: []string{"lo"},
		Short:   "Count non-test lines of code per file",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return metrics.RunLOC(os.Stdout, args, glob, exclude, marker)
		},
	}

	cmd.Flags().StringVarP(&glob, "glob", "g", "", "Include files matching glob")
	cmd.Flags().StringVarP(&exclude, "exclude", "e", "", "Exclude files with paths matching regex")
	cmd.Flags().StringVarP(&marker, "marker", "m", "", "Regex for start of test code")
	return cmd
}
