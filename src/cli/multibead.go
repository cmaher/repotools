package cli

import (
	"os"

	"repotools/src/beads"

	"github.com/spf13/cobra"
)

func newMultiBeadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "multi-bead id1 [id2 ...]",
		Short: "Show full details for multiple beads",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return beads.RunMultiBead(os.Stdout, args)
		},
	}
}
