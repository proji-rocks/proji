package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:   "add CLASS",
	Short: "add a new class",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("missing class name")
		}
		return class.AddClassCLI(args[0])
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}
