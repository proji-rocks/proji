package cmd

import (
	"github.com/spf13/cobra"
)

var classCmd = &cobra.Command{
	Use:   "class",
	Short: "add, remove and update your classes",
}

func init() {
	rootCmd.AddCommand(classCmd)
}
