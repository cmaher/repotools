package cli

import (
	"repotools/src/git"
	"repotools/src/runner"

	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log [base]",
		Short: "Commits since diverging from base branch",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			base := "master"
			if len(args) > 0 {
				base = args[0]
			}
			mb, err := git.MergeBase(base)
			if err != nil {
				return err
			}
			return runner.Exec([]string{"git", "log", "--oneline", mb + "..HEAD"})
		},
	}
}
