package proji

import (
	"context"
	"os"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
)

func projectCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clean",
		Short:                 "Auto-remove projects that have a dead path",
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cleanProjects(cmd.Context())
		},
	}

	return cmd
}

func doesPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func cleanProjects(ctx context.Context) error {
	logger := simplog.FromContext(ctx)

	// Get project manager from session
	logger.Debug("getting project manager from cli session")
	prma := cli.SessionFromContext(ctx).ProjectManager
	if prma == nil {
		return errors.New("no project manager available")
	}

	// Call the projects.
	logger.Debug("fetching project list")
	projects, err := prma.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, "fetch projects")
	}

	// Clean all projects in a pretty table.
	removeCounter := 0
	logger.Debug("cleaning projects")
	for _, project := range projects {
		if doesPathExist(project.Path) {
			continue // Skip if path exists
		}

		logger.Infof("Removing project %s (%q); previous location: %q", project.Name, project.ID, project.Path)
		if err = prma.Remove(ctx, project.ID); err != nil {
			return errors.Wrapf(err, "Failed to remove project %q", project.ID)
		}

		removeCounter++
	}

	if removeCounter > 0 {
		logger.Infof("Removed %d projects", removeCounter)
	} else {
		logger.Info("No projects to remove")
	}

	return nil
}
