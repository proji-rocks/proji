package cmd

import (
	"github.com/nikoksr/proji/internal/app/proji/global"

	"github.com/spf13/cobra"
)

var globalLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all globals",
	RunE: func(cmd *cobra.Command, args []string) error {
		return global.ListAll()
	},
}

func init() {
	globalCmd.AddCommand(globalLsCmd)
}
