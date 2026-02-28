package cli

import (
	"fmt"
	"os"

	"repotools/src/github"

	"github.com/spf13/cobra"
)

func newPRCmd() *cobra.Command {
	var only, exclude string

	cmd := &cobra.Command{
		Use:   "pr [number]",
		Short: "Show PR info, comments, reviews, checks, files, commits",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prArg := ""
			if len(args) > 0 {
				prArg = args[0]
			}

			sections := github.AllSections
			if only != "" {
				validated, err := github.ValidateSections(only)
				if err != nil {
					return err
				}
				sections = validated
			}
			if exclude != "" {
				sections = github.FilterSections(sections, "", exclude)
			}

			data, err := github.FetchPRData(prArg)
			if err != nil {
				return err
			}

			var reviewComments []github.ReviewComment
			for _, s := range sections {
				if s == "review-comments" {
					repo, err := github.GetRepoNWO()
					if err != nil {
						return err
					}
					reviewComments, _ = github.FetchReviewComments(repo, data.Number)
					break
				}
			}

			fmt.Fprintln(os.Stdout, github.RenderPR(*data, sections, reviewComments))
			return nil
		},
	}

	cmd.Flags().StringVar(&only, "only", "", "Comma-separated sections to show")
	cmd.Flags().StringVar(&exclude, "exclude", "", "Comma-separated sections to hide")
	return cmd
}
