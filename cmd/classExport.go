//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/storage/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exportAll, example bool
var destination string

var classExportCmd = &cobra.Command{
	Use:   "export LABEL [LABEL...]",
	Short: "Export one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if exportAll && example {
			return fmt.Errorf("the flags 'example' and 'all' cannot be passed at the same time")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Export an example class
		if example {
			file, err := exportExample(destination, projiEnv.ConfigFolderPath)
			if err != nil {
				fmt.Printf("> Export of example class failed: %v\n", err)
				return err
			}
			fmt.Printf("> Example class was successfully exported to %s\n", file)
			return nil
		}

		// Collect classes that will be exported
		var classes []*models.Class

		if exportAll {
			var err error
			classes, err = projiEnv.StorageService.LoadAllClasses()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing class label")
			}

			for _, label := range args {
				class, err := projiEnv.StorageService.LoadClass(label)
				if err != nil {
					return err
				}
				classes = append(classes, class)
			}
		}

		// Export the classes
		for _, class := range classes {
			if class.IsDefault {
				continue
			}
			fileOut, err := class.Export(destination)
			if err != nil {
				fmt.Printf("> Export of '%s' to file %s failed: %v\n", class.Label, fileOut, err)
				return err
			}
			fmt.Printf("> '%s' was successfully exported to file %s\n", class.Label, fileOut)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)
	classExportCmd.Flags().BoolVarP(&example, "example", "e", false, "Export an example class")
	classExportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all classes")

	classExportCmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination for the export")
	_ = classExportCmd.MarkFlagDirname("destination")
}

func exportExample(destination, confPath string) (string, error) {
	examplePath, ok := viper.Get("examples.path").(string)
	if !ok {
		return "", fmt.Errorf("could not read path of example config file")
	}

	examplePath = filepath.Join(confPath, examplePath)
	sourceFileStat, err := os.Stat(examplePath)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", examplePath)
	}

	src, err := os.Open(examplePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join(destination, "/proji-class-example.toml")
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return dstPath, err
}
