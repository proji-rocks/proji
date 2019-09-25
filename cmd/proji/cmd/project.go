package cmd

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "add, remove and update your projects",
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
