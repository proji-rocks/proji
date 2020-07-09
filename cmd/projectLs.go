package cmd

import (
	"io"
	"os"

	"github.com/nikoksr/proji/util"
	"github.com/pkg/errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type projectListCommand struct {
	cmd *cobra.Command
}

func newProjectListCommand() *projectListCommand {
	var cmd = &cobra.Command{
		Use:   "ls",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProjects(os.Stdout)
		},
	}
	return &projectListCommand{cmd: cmd}
}

func listProjects(out io.Writer) error {
	projects, err := activeSession.storageService.LoadProjects()
	if err != nil {
		return errors.Wrap(err, "failed to load all projects")
	}

	projectsTable := util.NewInfoTable(out)
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
