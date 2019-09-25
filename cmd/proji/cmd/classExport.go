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
	Short: "Export one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(exampleDest) > 0 {
			return storage.ExportExample(exampleDest)
		}

		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}
		for _, name := range args {
			var file string
			if file, err := ExportClass(name); err != nil {
				fmt.Printf("Export of class %s to file %s failed: %v\n", name, file, err)
				continue
			}
			fmt.Printf("Class %s was successfully exported to file %s.\n", name, file)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)
	classExportCmd.Flags().StringVarP(&exampleDest, "example", "e", "", "Export an example")
}

// ExportClass exports a class to a toml file.
// Returns the filename on success.
func ExportClass(name string) (string, error) {
	// Setup storage service
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return "", err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return "", err
	}
	defer s.Close()

	c, err := s.LoadClassByName(name)
	if err != nil {
		return "", err
	}
	return c.Export()
}
