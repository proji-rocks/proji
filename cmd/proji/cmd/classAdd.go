package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:   "add CLASS [CLASS...]",
	Short: "add new classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}

		for _, name := range args {
			if _, err := class.Add(name); err != nil {
				fmt.Printf("Failed adding class %s: %v\n", name, err)
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}
