/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package course

import (
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists the existing courses in the collection",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}
