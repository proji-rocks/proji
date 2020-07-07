//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of proji",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v" + session.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
