package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProjects()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listProjects() error {
	projects, err := projiEnv.Svc.LoadAllProjects()
	if err != nil {
		return err
	}

	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Install Path", "Class", "Status"})

	// Fill table
	for _, project := range projects {
		t.AppendRow([]interface{}{
			project.ID,
			project.Name,
			project.InstallPath,
			project.Class.Name,
			project.Status.Title,
		})
	}

	// Print the table
	t.Render()
	return nil
}
