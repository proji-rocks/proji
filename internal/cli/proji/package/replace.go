package pkg

import (
	"context"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
)

func newReplaceCommand() *cobra.Command {
	var forceReplacePackages bool

	cmd := &cobra.Command{
		Use:                   "replace [OPTIONS] LABEL PATH",
		Short:                 "Replace an already installed package with a new one",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return replacePackage(cmd.Context(), args[0], args[1])
		},
	}

	cmd.Flags().BoolVarP(&forceReplacePackages, "force", "f", false, "Don't ask for confirmation")

	return cmd
}

func replacePackage(ctx context.Context, label, config string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Load the package that should be replaced; this will also check if the package exists and helps us persist the
	// package's creation date.
	logger.Debugf("loading package %q", label)
	pkg, err := pama.GetByLabel(ctx, label)
	if err != nil {
		return errors.Wrapf(err, "load package %q", label)
	}

	// Before we remove, try to load the new package from the given config file. If this fails, we don't want to remove
	// the old package.
	logger.Debugf("loading new package from config file %q", config)
	newPkg, err := importPackage(ctx, config, true, nil)
	if err != nil {
		return errors.Wrapf(err, "load new package from config file %q", config)
	}

	// Remove the old package
	logger.Debugf("removing package %q", label)
	if err := pama.Remove(ctx, pkg.Label); err != nil {
		return errors.Wrapf(err, "remove package %q", label)
	}

	// Install the new package
	logger.Debugf("installing package %q", newPkg.Label)
	if err := pama.Store(ctx, newPkg); err != nil {
		return errors.Wrapf(err, "store package %q", newPkg.Label)
	}

	logger.Infof("Successfully replaced package %q", pkg.Label)

	return nil
}
