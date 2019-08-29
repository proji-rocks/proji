package cmd

import (
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "work on proji globals",
}

func init() {
	rootCmd.AddCommand(globalCmd)
}
