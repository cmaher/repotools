package cli

import (
	"os"

	"repotools/src/beads"

	"github.com/spf13/cobra"
)

func newBeadStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "bead-status",
		Aliases: []string{"bs"},
		Short:   "Print beads project status report to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return beads.RunBeadStatus(os.Stdout)
		},
	}
}
