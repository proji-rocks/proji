//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:    "completion SHELL",
	Short:  "Add shell completion",
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
