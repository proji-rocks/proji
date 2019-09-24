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
	Short: "import classes from config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing configfile name")
		}

		for _, config := range args {
			if err := ImportClass(config); err != nil {
				return err
			}
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
