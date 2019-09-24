package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var exampleDest string

var classExportCmd = &cobra.Command{
	Use:   "export CLASS [CLASS...]",
	Short: "export proji classes to config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(exampleDest) > 0 {
			return storage.ExportExample(exampleDest)
		}

		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}
		for _, name := range args {
			if err := ExportClass(name); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)
	classExportCmd.Flags().StringVarP(&exampleDest, "example", "e", "", "export example config")
}

// ExportClass exports a class to a toml file.
func ExportClass(name string) error {
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

	c, err := s.LoadClassByName(name)
	if err != nil {
		return err
	}
	return c.Export()
}
