package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:        "status",
	Short:      "Manage statuses",
	Deprecated: "support for project statuses will be dropped with v0.21.0",
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
