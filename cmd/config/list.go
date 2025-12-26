/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package config

import (
	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/service"
	"github.com/spf13/cobra"
)

func NewListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Command to list available configs",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = service.ListConfigs(cfg)
		},
	}
}
