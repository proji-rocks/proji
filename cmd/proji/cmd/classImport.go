package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing configfile name")
		}

		for _, config := range args {
			if err := ImportClass(config); err != nil {
				fmt.Printf("Import of file %s failed: %v\n", config, err)
				continue
			}
			fmt.Printf("File %s was successfully imported.\n", config)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)
}

// ImportClass imports a class from a config file.
func ImportClass(config string) error {
	// Import class data
	c, err := storage.NewClass("")
	if err != nil {
		return err
	}
	if err := c.ImportData(config); err != nil {
		return err
	}

	// Setup storage service
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.SaveClass(c)
}
