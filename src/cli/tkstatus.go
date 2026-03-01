package cli

import (
	"os"

	"repotools/src/tickets"

	"github.com/spf13/cobra"
)

func newTkStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tk-status",
		Aliases: []string{"ts"},
		Short:   "Print ticket project status report to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tickets.RunTicketStatus(os.Stdout)
		},
	}
}
