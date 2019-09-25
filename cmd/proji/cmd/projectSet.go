package cmd

import (
	"github.com/spf13/cobra"
)

var projectSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set project data",
}

func init() {
	rootCmd.AddCommand(projectSetCmd)
}
