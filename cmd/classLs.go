package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var classLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listClasses()
	},
}

func init() {
	classCmd.AddCommand(classLsCmd)
}

func listClasses() error {
	classes, err := projiEnv.Svc.LoadAllClasses()
	if err != nil {
		return err
	}

	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Label"})

	// Fill table
	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		t.AppendRow([]interface{}{class.Name, class.Label})
	}

	// Print the table
	t.Render()
	return nil
}
