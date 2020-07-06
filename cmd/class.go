//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"github.com/spf13/cobra"
)

var classCmd = &cobra.Command{
	Use:   "class",
	Short: "Manage classes",
}

func init() {
	rootCmd.AddCommand(classCmd)
}
