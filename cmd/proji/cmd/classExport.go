package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/spf13/cobra"
)

var exportExample bool

var classExportCmd = &cobra.Command{
	Use:   "export CLASS [CLASS...]",
	Short: "export proji classes to config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		numArgs := len(args)
		if exportExample {
			var destFolder string
			switch numArgs {
			case 0:
				destFolder = "."
			case 1:
				destFolder = args[0]
			default:
				return fmt.Errorf("invalid number of destination folders")
			}
			return class.ExportExample(destFolder)
		}

		if numArgs < 1 {
			return fmt.Errorf("missing class name")
		}
		for _, className := range args {
			if err := class.Export(className); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)

	// Flag to export an example config
	classExportCmd.Flags().BoolVarP(&exportExample, "example", "e", false, "export example config")
}
