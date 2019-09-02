package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:    "completion SHELL",
	Short:  "add shell completion",
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
