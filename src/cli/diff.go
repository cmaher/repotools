package cli

import (
	"repotools/src/git"
	"repotools/src/runner"

	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "diff [base] [flags...]",
		Short:              "Diff vs base branch",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			base := "master"
			extra := args
			if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
				base = args[0]
				extra = args[1:]
			}
			mb, err := git.MergeBase(base)
			if err != nil {
				return err
			}
			gitArgs := append([]string{"git", "diff", mb + "..HEAD"}, extra...)
			return runner.Exec(gitArgs)
		},
	}
	return cmd
}
