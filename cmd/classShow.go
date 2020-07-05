package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nikoksr/proji/pkg/util"

	"github.com/nikoksr/proji/pkg/storage/models"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var showAll bool

var classShowCmd = &cobra.Command{
	Use:   "show LABEL [LABEL...]",
	Short: "Show details about one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showAll {
			err := showAllClasses()
			if err != nil {
				fmt.Printf("> Showing of all classes failed: %v\n", err)
				return err
			}
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("missing class label")
		}

		for _, label := range args {
			err := showClass(nil, label)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
	classShowCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all classes")
}

func showClass(preloadedClass *models.Class, label string) error {
	var err error
	if preloadedClass == nil {
		preloadedClass, err = projiEnv.Svc.LoadClass(label)
		if err != nil {
			return nil
		}
	}
	output := os.Stdout
	showNameAndLabel(preloadedClass.Name, preloadedClass.Label)
	showTemplates(output, preloadedClass.Templates)
	showPlugins(output, preloadedClass.Plugins)
	return nil
}

func showAllClasses() error {
	classes, err := projiEnv.Svc.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		err = showClass(class, class.Label)
		if err != nil {
			fmt.Printf("failed printing table for class %s, %s\n", class.Name, err.Error())
		}
	}
	return nil
}

func showNameAndLabel(name, label string) {
	fmt.Printf("\nName:  %s\n", name)
	fmt.Printf("Label: %s\n\n", label)
}

func showTemplates(out io.Writer, templates []*models.Template) {
	templatesTable := util.NewInfoTable(out)
	templatesTable.SetTitle("TEMPLATES")
	templatesTable.AppendHeader(table.Row{"Destination", "Template Path", "Is File"})
	for _, template := range templates {
		templatesTable.AppendRow(table.Row{template.Destination, template.Path, template.IsFile})
	}
	templatesTable.Render()
}

func showPlugins(out io.Writer, plugins []*models.Plugin) {
	pluginsTable := util.NewInfoTable(out)
	pluginsTable.SetTitle("PLUGINS")
	pluginsTable.AppendHeader(table.Row{"Name", "Path", "Execution Number"})
	for _, plugin := range plugins {
		pluginsTable.AppendRow(table.Row{plugin.Name, plugin.Path, plugin.ExecNumber})
	}
	pluginsTable.Render()
}
