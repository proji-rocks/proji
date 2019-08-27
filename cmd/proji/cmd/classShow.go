package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classShowCmd = &cobra.Command{
	Use:   "show CLASS [CLASS...]",
	Short: "Show detailed class informations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}

		for _, className := range args {
			err := class.Show(className)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
}
