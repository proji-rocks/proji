package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classShowCmd = &cobra.Command{
	Use:   "show NAME [NAME...]",
	Short: "Show details about one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}

		for _, name := range args {
			if err := ShowClass(name); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
}

// ShowClass shows detailed information abour a given class
func ShowClass(name string) error {
	// Setup storage service
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	c, err := s.LoadClassByName(name)
	if err != nil {
		return err
	}

	// fmt.Println(helper.ProjectHeader(c.Name))
	showLabels(c)
	showFolders(c)
	showFiles(c)
	showScripts(c)
	return nil
}

// showLabels shows all labels of a given class
func showLabels(class *storage.Class) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Label"})

	for _, label := range class.Labels {
		t.AppendRow([]interface{}{label})
	}
	// Print the table
	t.Render()
	fmt.Println()
}

// showFolders shows all folders of a given class
func showFolders(class *storage.Class) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Folder", "Template"})

	// Fill table
	for folder, template := range class.Folders {
		t.AppendRow([]interface{}{folder, template})
	}

	// Print the table
	t.Render()
	fmt.Println()
}

// showFiles shows all files of a given class
func showFiles(class *storage.Class) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"File", "Template"})

	// Fill table
	for folder, template := range class.Files {
		t.AppendRow([]interface{}{folder, template})
	}

	// Print the table
	t.Render()
	fmt.Println()
}

// showScripts shows all scripts of a given class
func showScripts(class *storage.Class) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Script", "As sudo"})

	// Fill table
	for script, runAsSudo := range class.Scripts {
		t.AppendRow([]interface{}{script, runAsSudo})
	}

	// Print the table
	t.Render()
	fmt.Println()
}
