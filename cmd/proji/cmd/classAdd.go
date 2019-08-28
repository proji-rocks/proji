package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:   "add CLASS",
	Short: "Add a new class",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Missing class name")
		}

		err := class.AddClassCLI(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}
