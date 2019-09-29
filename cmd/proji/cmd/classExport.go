package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exampleDest string

var classExportCmd = &cobra.Command{
	Use:   "export LABEL [LABEL...]",
	Short: "Export one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(exampleDest) > 0 {
			return ExportExample(exampleDest)
		}

		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}
		for _, label := range args {
			file, err := ExportClass(label)
			if err != nil {
				fmt.Printf("Export of '%s' to file %s failed: %v\n", label, file, err)
				continue
			}
			fmt.Printf("'%s' was successfully exported to file %s.\n", label, file)
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
func ExportClass(label string) (string, error) {
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

	classID, err := s.LoadClassIDByLabel(label)
	if err != nil {
		return "", err
	}
	class, err := s.LoadClass(classID)
	if err != nil {
		return "", err
	}
	return class.Export()
}

// ExportExample exports an example class config
func ExportExample(destFolder string) error {

	exampleDir, ok := viper.Get("examples.location").(string)
	if !ok {
		return fmt.Errorf("Could not read example file location from config file")
	}
	exampleFile, ok := viper.Get("examples.class").(string)
	if !ok {
		return fmt.Errorf("Could not read example file name from config file")
	}

	exampleFile = helper.GetConfigDir() + exampleDir + exampleFile
	sourceFileStat, err := os.Stat(exampleFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", exampleFile)
	}

	source, err := os.Open(exampleFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destFolder + "/proji-class.toml")
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
