package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/messages"

	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"
	"github.com/nikoksr/proji/util"

	"github.com/spf13/cobra"
)

type packageRemoveCommand struct {
	cmd *cobra.Command
}

func newPackageRemoveCommand() *packageRemoveCommand {
	var removeAllPackages, forceRemovePackages bool

	var cmd = &cobra.Command{
		Use:   "rm LABEL [LABEL...]",
		Short: "Remove one or more packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Collect packages that will be removed
			var packages []*models.Package

			if removeAllPackages {
				var err error
				packages, err = activeSession.storageService.LoadPackages()
				if err != nil {
					return errors.Wrap(err, "failed to load all packages")
				}
			} else {
				if len(args) < 1 {
					return fmt.Errorf("missing package label")
				}

				for _, label := range args {
					pkg, err := activeSession.storageService.LoadPackage(label)
					if err != nil {
						messages.Warningf("failed to load package, %s", err.Error())
						continue
					}
					packages = append(packages, pkg)
				}
			}

			// Remove the packages
			for _, pkg := range packages {
				// Skip default packages
				if pkg.IsDefault {
					continue
				}
				// Ask for confirmation if force flag was not passed
				if !forceRemovePackages {
					if !util.WantTo(
						fmt.Sprintf("Do you really want to remove package '%s (%s)'?", pkg.Name, pkg.Label),
					) {
						continue
					}
				}
				err := activeSession.storageService.RemovePackage(pkg.Label)
				if err != nil {
					messages.Warningf("failed to remove package %s, %s", pkg.Label, err.Error())
				} else {
					messages.Successf("successfully remove package %s", pkg.Label)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&removeAllPackages, "all", "a", false, "Remove all packages")
	cmd.Flags().BoolVarP(&forceRemovePackages, "force", "f", false, "Don't ask for confirmation")
	return &packageRemoveCommand{cmd: cmd}
}
