package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classRmCmd = &cobra.Command{
	Use:   "rm CLASS [CLASS...]",
	Short: "Remove existing classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}

		for _, className := range args {
			err := class.RemoveClass(className)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
}
