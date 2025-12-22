/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package cmd

import (
	"os"

	"github.com/dhanush-cache/course-flow/cmd/course"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "course-flow",
	Short: "A powerful CLI to manage and organize your educational courses.",
	Long: `Course-Flow is a command-line interface (CLI) tool designed to help you efficiently manage your collection of educational courses, whether they are from online platforms, universities, or local resources.
Use 'course-flow [command] --help' for more information about a specific command.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(course.Cmd)
}
