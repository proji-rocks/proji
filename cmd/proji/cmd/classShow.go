package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classShowCmd = &cobra.Command{
	Use:   "show CLASS [CLASS...]",
	Short: "show detailed class informations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}

		for _, className := range args {
			if err := class.Show(className); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
}
