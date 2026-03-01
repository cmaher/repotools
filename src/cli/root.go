package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var directory string

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repotools",
		Short: "Repo helper toolkit: git, GitHub, and filesystem operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if directory != "" {
				return os.Chdir(directory)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	cmd.PersistentFlags().StringVarP(&directory, "directory", "C", "", "Change to DIR before doing anything")

	cmd.AddCommand(
		newStatusCmd(),
		newLogCmd(),
		newDiffCmd(),
		newLsCmd(),
		newPRCmd(),
		newReadCmd(),
		newMultiLSCmd(),
		newMultiFindCmd(),
		newLocCmd(),
		newFnSpansCmd(),
		newTkStatusCmd(),
	)

	return cmd
}
