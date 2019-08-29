package cmd

import (
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Work on proji globals",
}

func init() {
	rootCmd.AddCommand(globalCmd)
}
