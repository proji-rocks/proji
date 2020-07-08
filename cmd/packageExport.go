//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/messages"

	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exportAll, example bool
var destination string

var packageExportCmd = &cobra.Command{
	Use:   "export LABEL [LABEL...]",
	Short: "ExportConfig one or more packages",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if exportAll && example {
			return fmt.Errorf("the flags 'example' and 'all' cannot be passed at the same time")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// ExportConfig an example package
		if example {
			file, err := exportExample(destination, session.Config.BasePath)
			if err != nil {
				return errors.Wrap(err, "failed to export example package")
			}
			messages.Success("successfully exported example package to %s", file)
			return nil
		}

		// Collect packages that will be exported
		var packages []*models.Package

		if exportAll {
			var err error
			packages, err = session.StorageService.LoadPackages()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing package label")
			}

			for _, label := range args {
				pkg, err := session.StorageService.LoadPackage(label)
				if err != nil {
					return err
				}
				packages = append(packages, pkg)
			}
		}

		// ExportConfig the packages
		for _, pkg := range packages {
			if pkg.IsDefault {
				continue
			}
			fileOut, err := pkg.ExportConfig(destination)
			if err != nil {
				messages.Warning("failed to export package %s to %s, %s", pkg.Label, fileOut, err.Error())
			} else {
				messages.Success("successfully exported package %s to %s", pkg.Label, fileOut)
			}
		}
		return nil
	},
}

func init() {
	packageCmd.AddCommand(packageExportCmd)
	packageExportCmd.Flags().BoolVarP(&example, "example", "e", false, "ExportConfig an example package")
	packageExportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "ExportConfig all packages")

	packageExportCmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination for the export")
	_ = packageExportCmd.MarkFlagDirname("destination")
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

	dstPath := filepath.Join(destination, "/proji-package-example.toml")
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return dstPath, err
}
