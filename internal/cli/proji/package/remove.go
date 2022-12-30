package pkg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/nikoksr/simplog"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
)

func newRemoveCommand() *cobra.Command {
	var forceRemovePackages bool

	cmd := &cobra.Command{
		Use:                   "rm [OPTIONS] LABEL [LABEL...]",
		Short:                 "Remove installed packages",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return removePackages(cmd.Context(), args...)
		},
	}

	cmd.Flags().BoolVarP(&forceRemovePackages, "force", "f", false, "Don't ask for confirmation")

	return cmd
}

func removePackages(ctx context.Context, labels ...string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Removing packages
	logger.Debugf("removing %d packages", len(labels))
	for _, label := range labels {
		logger.Debugf("removing package %q", label)
		if err := pama.Remove(ctx, label); err != nil {
			return errors.Wrapf(err, "remove package %q", label)
		}
	}

	return nil
}
