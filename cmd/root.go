/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package cmd

import (
	"log"
	"os"

	"github.com/dhanush-cache/course-flow/cmd/course"
	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/spf13/cobra"
)

var cfg *config.Config

func NewRootCommand() *cobra.Command {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	rootCmd := &cobra.Command{
		Use:   "course-flow",
		Short: "A powerful CLI to manage and organize your educational courses.",
		Long: `Course-Flow is a command-line interface (CLI) tool designed to help you efficiently manage your collection of educational courses, whether they are from online platforms, universities, or local resources.
Use 'course-flow [command] --help' for more information about a specific command.`,
	}

	rootCmd.AddCommand(course.NewCommand(cfg))

	return rootCmd
}

func Execute() {
	err := NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
