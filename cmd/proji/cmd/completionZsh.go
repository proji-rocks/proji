package cmd

import (
	"github.com/spf13/cobra"
)

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "zsh completion",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenZshCompletionFile("proji-zsh-completion")
	},
}

func init() {
	completionCmd.AddCommand(completionZshCmd)
}
