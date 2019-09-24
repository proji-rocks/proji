package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classShowCmd = &cobra.Command{
	Use:   "show CLASS [CLASS...]",
	Short: "show detailed class informations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
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

	fmt.Println(helper.ProjectHeader(c.Name))
	showLabels(c)
	showFolders(c)
	showFiles(c)
	showScripts(c)
	return nil
}

// showLabels shows all labels of a given class
func showLabels(class *storage.Class) {
	fmt.Println("Labels:")

	for _, label := range class.Labels {
		fmt.Println(" " + label)
	}
	fmt.Println()
}

// showFolders shows all folders of a given class
func showFolders(class *storage.Class) {
	fmt.Println("Folders:")

	for target, template := range class.Folders {
		fmt.Println(" " + target + " : " + template)
	}
	fmt.Println()
}

// showFiles shows all files of a given class
func showFiles(class *storage.Class) {
	fmt.Println("Files:")

	for target, template := range class.Files {
		fmt.Println(" " + target + " : " + template)
	}
	fmt.Println()
}

// showScripts shows all scripts of a given class
func showScripts(class *storage.Class) {
	fmt.Println("Scripts:")

	for script, runAsSudo := range class.Scripts {
		sudo := ""
		if runAsSudo {
			sudo = "sudo "
		}
		fmt.Println(" " + sudo + script)
	}
	fmt.Println()
}
