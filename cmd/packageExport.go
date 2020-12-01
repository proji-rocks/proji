package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/internal/static"
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type packageExportCommand struct {
	cmd *cobra.Command
}

func newPackageExportCommand() *packageExportCommand {
	var exportAll, template, asJson bool
	var destination string

	cmd := &cobra.Command{
		Use:     "export LABEL [LABEL...]",
		Short:   "Export one or more packages",
		Aliases: []string{"e"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if exportAll && template {
				return fmt.Errorf("the flags 'template' and 'all' cannot be used at the same time")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Export an example package
			if template {
				file, err := exportExample(destination)
				if err != nil {
					return errors.Wrap(err, "failed to export config template")
				}
				message.Successf("successfully exported config template to %s", file)
				return nil
			}

			// Collect packages that will be exported
			var packages []*domain.Package
			var err error

			if exportAll {
				packages, err = session.packageService.LoadPackageList(true)
				if err != nil {
					return err
				}
			} else {
				if len(args) < 1 {
					return fmt.Errorf("missing package label")
				}
				packages, err = session.packageService.LoadPackageList(true, args...)
				if err != nil {
					return err
				}
			}

			// Export the packages
			for _, pkg := range packages {
				exportedTo, err := session.packageService.ExportPackageToConfig(*pkg, ".", asJson)
				if err != nil {
					message.Warningf("failed to export package %s to %s, %v", pkg.Label, exportedTo, err)
				} else {
					message.Successf("successfully exported package %s to %s", pkg.Label, exportedTo)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&template, "example", "e", false, "Export a package config example")
	cmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all packages")
	cmd.Flags().BoolVarP(&asJson, "json", "j", false, "Export to a Json file")
	cmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination for the export")

	_ = cmd.MarkFlagDirname("destination")

	return &packageExportCommand{cmd: cmd}
}

func exportExample(destination string) (string, error) {
	destination = filepath.Join(destination, "proji-package-example.toml")
	file, err := os.Create(destination)
	if err != nil {
		return "", errors.Wrap(err, "create config template")
	}
	defer file.Close()
	_, err = file.WriteString(static.PackageConfigTemplate)
	return destination, err
}
