/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package course

import (
	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/spf13/cobra"
)

func NewAddCommand(cfg *config.Config) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new course to your collection",
		Long:  `Download and organize a new course into your course collection.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			files, _ := cmd.Flags().GetStringSlice("files")
			_ = key
			_ = files
			_ = cfg
		},
	}

	addCmd.Flags().StringSliceP("files", "f", []string{}, "Path to the course zip file(s)")
	_ = addCmd.MarkFlagRequired("files")

	return addCmd
}
