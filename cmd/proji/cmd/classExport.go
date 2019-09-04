package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var exampleDest string

var classExportCmd = &cobra.Command{
	Use:   "export CLASS [CLASS...]",
	Short: "export proji classes to config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(exampleDest) > 0 {
			return class.ExportExample(exampleDest)
		}

		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}
		for _, name := range args {
			c := class.New(name)
			if err := c.Export(); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)
	classExportCmd.Flags().StringVarP(&exampleDest, "example", "e", ".", "export example config")
}
