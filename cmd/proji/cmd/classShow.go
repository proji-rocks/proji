package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classShowCmd = &cobra.Command{
	Use:   "show LABEL [LABEL...]",
	Short: "Show details about one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}

		for _, name := range args {
			if err := showClass(name); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
}

func showClass(label string) error {
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

	classID, err := s.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}
	class, err := s.LoadClass(classID)
	if err != nil {
		return nil
	}

	// fmt.Println(helper.ProjectHeader(c.Name))
	showInfo(class.Name, class.Label)
	showFolders(class.Folders)
	showFiles(class.Files)
	showScripts(class.Scripts)
	return nil
}

func showInfo(name, label string) {
	fmt.Println("Name: " + name)
	fmt.Println("Label: " + label)
	fmt.Println()
}

func showFolders(folders map[string]string) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Folder", "Template"})

	// Fill table
	for folder, template := range folders {
		t.AppendRow([]interface{}{folder, template})
	}

	// Print the table
	t.Render()
}

func showFiles(files map[string]string) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"File", "Template"})

	// Fill table
	for folder, template := range files {
		t.AppendRow([]interface{}{folder, template})
	}

	// Print the table
	t.Render()
}

func showScripts(scripts map[string]bool) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Script", "As sudo"})

	// Fill table
	for script, runAsSudo := range scripts {
		t.AppendRow([]interface{}{script, runAsSudo})
	}

	// Print the table
	t.Render()
}
