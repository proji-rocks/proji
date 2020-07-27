package cmd

import (
	"os"

	"github.com/nikoksr/proji/internal/util"
	"github.com/pkg/errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type projectListCommand struct {
	cmd *cobra.Command
}

func newProjectListCommand() *projectListCommand {
	cmd := &cobra.Command{
		Use:                   "ls",
		Short:                 "List projects",
		Aliases:               []string{"l"},
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProjects()
		},
	}
	return &projectListCommand{cmd: cmd}
}

func listProjects() error {
	projects, err := session.projectService.LoadProjectList()
	if err != nil {
		return errors.Wrap(err, "failed to load all projects")
	}

	projectsTable := util.NewInfoTable(os.Stdout)
	projectsTable.AppendHeader(table.Row{"Name", "Install Path", "Package"})

	for _, project := range projects {
		projectsTable.AppendRow(table.Row{
			project.Name,
			project.Path,
			project.Package.Name,
		})
	}

	projectsTable.Render()
	return nil
}
