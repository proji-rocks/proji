//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"io"
	"os"

	"github.com/nikoksr/proji/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	PreRun: func(cmd *cobra.Command, args []string) {
		setMaxColumnWidth()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProjects(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listProjects(out io.Writer) error {
	projects, err := session.StorageService.LoadProjects()
	if err != nil {
		return err
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
