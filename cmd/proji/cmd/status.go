package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Manage statuses",
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
