package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/internal/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type packageRemoveCommand struct {
	cmd *cobra.Command
}

func newPackageRemoveCommand() *packageRemoveCommand {
	var removeAllPackages, forceRemovePackages bool

	cmd := &cobra.Command{
		Use:   "rm LABEL [LABEL...]",
		Short: "Remove one or more packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Collect packages that will be removed
			var err error
			var labels []string

			// Determine labels of packages that will be removed
			if removeAllPackages {
				labels, err = getAllPackageLabels()
				if err != nil {
					return err
				}
			} else {
				if len(args) < 1 {
					return fmt.Errorf("missing package label")
				}
				labels = args
			}

			// Remove the packages
			for _, label := range labels {
				// Ask for confirmation if force flag was not passed
				if !forceRemovePackages {
					if !util.WantTo(
						fmt.Sprintf("Do you really want to remove package %s?", label),
					) {
						continue
					}
				}
				err := session.packageService.RemovePackage(label)
				if err != nil {
					message.Warningf("failed to remove package %s, %v", label, err)
				} else {
					message.Successf("successfully removed package %s", label)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&removeAllPackages, "all", "a", false, "Remove all packages")
	cmd.Flags().BoolVarP(&forceRemovePackages, "force", "f", false, "Don't ask for confirmation")
	return &packageRemoveCommand{cmd: cmd}
}

func getAllPackageLabels() ([]string, error) {
	pkgs, err := session.packageService.LoadPackageList(false)
	if err != nil {
		return nil, errors.Wrap(err, "load packages")
	}

	labels := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		labels = append(labels, pkg.Label)
	}
	return labels, nil
}
