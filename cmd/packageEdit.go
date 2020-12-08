package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/internal/message"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type packageEditCommand struct {
	cmd *cobra.Command
}

func newPackageEditCommand() *packageEditCommand {
	var packageLabel string

	cmd := &cobra.Command{
		Use:     "edit LABEL",
		Short:   "Edit a package config",
		Aliases: []string{"e"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("no package label given")
			} else if len(args) > 1 {
				message.Warningf("more than one package label given. Proceeding with first label and ignoring the rest.")
			}

			packageLabel = args[0]
			if len(packageLabel) < 1 {
				return fmt.Errorf("given label is empty")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := editPackageConfig(packageLabel)
			if err != nil {
				message.Warningf("failed to edit package %s: %v", packageLabel, err)
			} else {
				message.Successf("successfully edited package %s", packageLabel)
			}
		},
	}

	return &packageEditCommand{cmd: cmd}
}

func editPackageConfig(packageLabel string) (err error) {
	// Try to load the package; implicitly checks if a package is even associated to the given label.
	pkg, err := session.packageService.LoadPackage(true, packageLabel)
	if err != nil {
		return fmt.Errorf("load package: %v", err)
	}

	// Export package to temporary config file.
	configFile, err := session.packageService.ExportPackageToTemporaryConfig(*pkg)
	if err != nil {
		return fmt.Errorf("export package to temporary config file: %v", err)
	}
	defer func() {
		ferr := os.Remove(configFile)
		if ferr == nil {
			return
		}
		if err != nil {
			err = errors.Wrap(err, ferr.Error())
		} else {
			err = ferr
		}
	}()

	// Open config in system's default text editor and wait for it to start.
	err = config.OpenInEditor(configFile)
	if err != nil {
		return fmt.Errorf("open package config in text editor: %v", err)
	}

	fmt.Print("Press 'Enter' after you have saved the changes to the config file...")
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("reading 'Enter' key press: %v", err)
	}

	// Try to import the new/edited config.
	oldLabel := pkg.Label
	pkg, err = session.packageService.ImportPackageFromConfig(configFile)
	if err != nil {
		return fmt.Errorf("import edited config: %v", err)
	}

	// Replace package with new config.
	err = session.packageService.RemovePackage(oldLabel)
	if err != nil {
		return fmt.Errorf("remove package: %v", err)
	}

	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return fmt.Errorf("store edited config: %v", err)
	}

	return nil
}
