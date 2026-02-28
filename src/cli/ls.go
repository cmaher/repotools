package cli

import (
	"repotools/src/git"
	"repotools/src/runner"

	"github.com/spf13/cobra"
)

func newLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "ls [base] [-- path...]",
		Short:              "List files at merge base",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			base := "master"
			extra := args
			if len(args) > 0 && args[0] != "" && args[0] != "--" && args[0][0] != '-' {
				base = args[0]
				extra = args[1:]
			}
			mb, err := git.MergeBase(base)
			if err != nil {
				return err
			}
			gitArgs := append([]string{"git", "ls-tree", "--name-only", mb}, extra...)
			return runner.Exec(gitArgs)
		},
	}
	return cmd
}
