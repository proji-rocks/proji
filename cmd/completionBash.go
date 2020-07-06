//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Bash completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
}
