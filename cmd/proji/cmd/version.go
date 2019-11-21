package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of proji",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v0.18.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
