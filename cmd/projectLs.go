package cmd

import (
	"io"
	"os"

	"github.com/nikoksr/proji/pkg/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProjects(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listProjects(out io.Writer) error {
	projects, err := projiEnv.StorageService.LoadAllProjects()
	if err != nil {
		return err
	}

	projectsTable := util.NewInfoTable(out)
	projectsTable.AppendHeader(table.Row{"ID", "Name", "Install Path", "Class"})

	for _, project := range projects {
		projectsTable.AppendRow(table.Row{
			project.ID,
			project.Name,
			project.Path,
			project.Class.Name,
		})
	}

	projectsTable.Render()
	return nil
}
