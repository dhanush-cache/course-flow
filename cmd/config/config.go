/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package config

import (
	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Command to manage configuration settings",
		Long:  `List, view, and modify configuration settings for Course-Flow CLI tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	configCmd.AddCommand(NewListCommand(cfg))
	return configCmd
}
