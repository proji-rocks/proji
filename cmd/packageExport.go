package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type packageExportCommand struct {
	cmd *cobra.Command
}

func newPackageExportCommand() *packageExportCommand {
	var exportAll, example bool
	var destination string

	var cmd = &cobra.Command{
		Use:   "export LABEL [LABEL...]",
		Short: "Export one or more packages",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if exportAll && example {
				return fmt.Errorf("the flags 'example' and 'all' cannot be passed at the same time")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Export an example package
			if example {
				file, err := exportExample(destination, session.config.BasePath)
				if err != nil {
					return errors.Wrap(err, "failed to export example package")
				}
				message.Successf("successfully exported example package to %s", file)
				return nil
			}

			// Collect packages that will be exported
			var packages []*domain.Package
			var err error

			if exportAll {
				packages, err = session.packageService.LoadPackageList()
				if err != nil {
					return err
				}
			} else {
				if len(args) < 1 {
					return fmt.Errorf("missing package label")
				}
				packages, err = session.packageService.LoadPackageList(args...)
				if err != nil {
					return err
				}
			}

			// Export the packages
			for _, pkg := range packages {
				if pkg.IsDefault {
					continue
				}

				exportedTo, err := session.packageService.ExportPackageToConfig(*pkg, ".")
				if err != nil {
					message.Warningf("failed to export package %s to %s, %v", pkg.Label, exportedTo, err)
				} else {
					message.Successf("successfully exported package %s to %s", pkg.Label, exportedTo)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&example, "example", "e", false, "Export an example package")
	cmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all packages")
	cmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination for the export")
	_ = cmd.MarkFlagDirname("destination")

	return &packageExportCommand{cmd: cmd}
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
