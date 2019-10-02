package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing config file")
		}

		for _, config := range args {
			if err := importClass(config, projiEnv.Svc); err != nil {
				fmt.Printf("Import of '%s' failed: %v\n", config, err)
				continue
			}
			fmt.Printf("'%s' was successfully imported.\n", config)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)
}

func importClass(config string, svc storage.Service) error {
	// Import class data
	c, err := storage.NewClass("", "")
	if err != nil {
		return err
	}
	if err := c.ImportData(config); err != nil {
		return err
	}
	return svc.SaveClass(c)
}
