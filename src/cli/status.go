package cli

import (
	"os"

	"repotools/src/git"

	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current branch and working tree status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return git.Status(os.Stdout)
		},
	}
}
