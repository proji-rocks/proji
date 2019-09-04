package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var classRmCmd = &cobra.Command{
	Use:   "rm CLASS [CLASS...]",
	Short: "remove existing classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}

		for _, name := range args {
			c := class.New(name)
			if err := c.Remove(); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
}
