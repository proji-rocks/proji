//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/messages"

	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"
	"github.com/nikoksr/proji/util"

	"github.com/spf13/cobra"
)

var removeAllPackages, forceRemovePackages bool

var packageRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Collect packages that will be removed
		var packages []*models.Package

		if removeAllPackages {
			var err error
			packages, err = session.StorageService.LoadPackages()
			if err != nil {
				return errors.Wrap(err, "failed to load all packages")
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing package label")
			}

			for _, label := range args {
				pkg, err := session.StorageService.LoadPackage(label)
				if err != nil {
					messages.Warning("failed to load package, %s", err.Error())
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
			err := session.StorageService.RemovePackage(pkg.Label)
			if err != nil {
				messages.Warning("failed to remove package %s, %s", pkg.Label, err.Error())
			} else {
				messages.Success("successfully remove package %s", pkg.Label)
			}
		}
		return nil
	},
}

func init() {
	packageCmd.AddCommand(packageRmCmd)
	packageRmCmd.Flags().BoolVarP(&removeAllPackages, "all", "a", false, "Remove all packages")
	packageRmCmd.Flags().BoolVarP(&forceRemovePackages, "force", "f", false, "Don't ask for confirmation")
}
