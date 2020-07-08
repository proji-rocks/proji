//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Manage packages",
}

func init() {
	rootCmd.AddCommand(packageCmd)
}
