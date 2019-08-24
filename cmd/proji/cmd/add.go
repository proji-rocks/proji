package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add COMMAND ARGS",
	Short: "Add new instances of a specified type.",
}

func init() {
	rootCmd.AddCommand(addCmd)
}
