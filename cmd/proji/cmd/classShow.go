package cmd

import (
	"fmt"
	"os"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/jedib0t/go-pretty/table"
	"github.com/nikoksr/proji/pkg/proji/storage"
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
			if err := showClass(name, projiEnv.Svc); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
}

func showClass(label string, svc storage.Service) error {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}
	class, err := svc.LoadClass(classID)
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
	fmt.Println("\nName: " + name)
	fmt.Println("Label: " + label)
	fmt.Println()
}

func showFolders(folders []*item.Folder) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Folder", "Template"})

	// Fill table
	for _, folder := range folders {
		t.AppendRow([]interface{}{folder.Destination, folder.Template})
	}

	// Print the table
	t.Render()
}

func showFiles(files []*item.File) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"File", "Template"})

	// Fill table
	for _, file := range files {
		t.AppendRow([]interface{}{file.Destination, file.Template})
	}

	// Print the table
	t.Render()
}

func showScripts(scripts []*item.Script) {
	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Script", "As sudo"})

	// Fill table
	for _, script := range scripts {
		t.AppendRow([]interface{}{script.ExecNumber, script.Name, script.RunAsSudo})
	}

	// Print the table
	t.Render()
}
