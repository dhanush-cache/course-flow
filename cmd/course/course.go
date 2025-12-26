/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package course

import (
	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "course",
		Short: "Commands for working with the collection of courses",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(NewAddCommand(cfg))
	cmd.AddCommand(NewListCommand())

	return cmd
}
