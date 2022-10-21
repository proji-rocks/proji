package proji

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/internal/text"
	"github.com/nikoksr/proji/pkg/logging"
)

func projectListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ls",
		Short:                 "List previously created projects",
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return listProjects(cmd.Context())
		},
	}

	return cmd
}

func listProjects(ctx context.Context) error {
	logger := logging.FromContext(ctx)

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

	// Exit when no projects are installed.
	if len(projects) == 0 {
		return nil
	}

	// List all projects in a pretty table.
	logger.Debug("listing projects")

	fmt.Println()
	table := text.NewTablePrinter()
	table.AddHeaderColumns("#", "ID", "Name", "Package", "Path", "Created at")

	for idx, project := range projects {
		table.AddRow(idx+1, project.ID, project.Name, project.Package, project.Path, project.CreatedAt)
	}

	err = table.Render()
	if err != nil {
		return errors.Wrap(err, "render table")
	}

	return nil
}
