/*
Copyright © 2025 Dhanush Shetty dhanushshettycache@outlook.com
*/

package course

import (
	"fmt"
	"log"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/service"
	"github.com/spf13/cobra"
)

func NewAddCommand(cfg *config.Config) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new course to your collection",
		Long:  `Download and organize a new course into your course collection.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			filesSet := cmd.Flags().Changed("files")
			directorySet := cmd.Flags().Changed("directory")
			if !filesSet && !directorySet {
				return fmt.Errorf("exactly one of --files or --directory must be provided")
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			courseID := args[0]
			files, _ := cmd.Flags().GetStringSlice("files")
			directory, _ := cmd.Flags().GetString("directory")
			options := service.Options{
				CourseID: courseID,
			}
			if cmd.Flags().Changed("files") {
				options.ZipFiles = files
			} else if cmd.Flags().Changed("directory") {
				options.Directory = directory
			}
			err := service.AddCourse(options, cfg)
			if err != nil {
				log.Fatalf("Error adding course: %v", err)
			}
		},
	}

	addCmd.Flags().StringSliceP("files", "f", []string{}, "Path to the course zip file(s)")
	addCmd.Flags().StringP("directory", "d", "", "Path to the course source directory (extracted)")
	addCmd.MarkFlagsMutuallyExclusive("files", "directory")

	return addCmd
}
