package cmd

import (
	"github.com/spf13/cobra"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "bash completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletionFile("proji-bash-completion")
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
}
