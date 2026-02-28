package cli

import (
	"os"

	"repotools/src/metrics"

	"github.com/spf13/cobra"
)

func newFnSpansCmd() *cobra.Command {
	var glob, excludePath, pattern, after, include, exclude string

	cmd := &cobra.Command{
		Use:   "fn-spans paths...",
		Short: "Show function line ranges in source files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return metrics.RunFnSpans(os.Stdout, args, glob, excludePath, pattern, after, include, exclude)
		},
	}

	cmd.Flags().StringVarP(&glob, "glob", "g", "", "Include files matching glob")
	cmd.Flags().StringVarP(&excludePath, "exclude-path", "E", "", "Exclude files with paths matching regex")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Regex for function defs (group 1 = name)")
	cmd.Flags().StringVarP(&after, "after", "a", "", "Only scan after first line matching this")
	cmd.Flags().StringVarP(&include, "include", "i", "", "Only include functions matching regex")
	cmd.Flags().StringVarP(&exclude, "exclude", "x", "", "Exclude functions matching regex")
	return cmd
}
