/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package course

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "course",
	Short: "Commands for working with the collection of courses",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(listCmd)
}
