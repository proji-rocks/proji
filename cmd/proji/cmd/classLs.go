package cmd

import (
	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var classLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List existing classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := class.ListAll()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classLsCmd)
}
