package proji

import (
	"context"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
)

func projectRemoveCommand() *cobra.Command {
	var forceRemoveProjects bool

	cmd := &cobra.Command{
		Use:                   "rm [OPTIONS] ID [ID...]",
		Short:                 "Remove tracked projects",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return removeProjects(cmd.Context(), args...)
		},
	}

	cmd.Flags().BoolVarP(&forceRemoveProjects, "force", "f", false, "Don't ask for confirmation")

	return cmd
}

func removeProjects(ctx context.Context, ids ...string) error {
	logger := simplog.FromContext(ctx)

	// Get project manager from session
	logger.Debug("getting project manager from cli session")
	prma := cli.SessionFromContext(ctx).ProjectManager
	if prma == nil {
		return errors.New("no project manager found")
	}

	// Removing projects
	logger.Debugf("removing %d projects", len(ids))
	for _, id := range ids {
		logger.Debugf("removing project %q", id)
		if err := prma.Remove(ctx, id); err != nil {
			logger.Warnf("Failed to remove project %q: %v", id, err)
		} else {
			logger.Infof("Successfully removed project %q", id)
		}
	}

	return nil
}
