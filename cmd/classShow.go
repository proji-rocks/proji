//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/nikoksr/proji/util"

	"github.com/nikoksr/proji/storage/models"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var showAll bool

var classShowCmd = &cobra.Command{
	Use:   "show LABEL [LABEL...]",
	Short: "Show details about one or more classes",
	PreRun: func(cmd *cobra.Command, args []string) {
		setMaxColumnWidth()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !showAll && len(args) < 1 {
			return fmt.Errorf("missing class label")
		}

		var labels []string
		if !showAll {
			labels = args
		}
		return showClasses(labels...)
	},
}

func init() {
	classCmd.AddCommand(classShowCmd)
	classShowCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all classes")
}

func showClass(preloadedClass *models.Class, label string) error {
	var err error
	if preloadedClass == nil {
		preloadedClass, err = session.StorageService.LoadClass(label)
		if err != nil {
			return err
		}
	}
	output := os.Stdout
	showBasicInfo(preloadedClass.Name, preloadedClass.Label, preloadedClass.Description)
	showTemplates(output, preloadedClass.Templates)
	showPlugins(output, preloadedClass.Plugins)
	return nil
}

func showClasses(labels ...string) error {
	classes, err := session.StorageService.LoadClasses(labels...)
	if err != nil {
		return err
	}
	for _, class := range classes {
		err = showClass(class, class.Label)
		if err != nil {
			log.Printf("failed show class with label '%s', %s", class.Label, err.Error())
		}
	}
	return nil
}

func showBasicInfo(name, label, description string) {
	fmt.Printf("\nName:  %s\n", name)
	fmt.Printf("Label: %s\n", label)
	fmt.Printf("Description: %s\n\n", text.WrapSoft(description, maxColumnWidth))
}

func showTemplates(out io.Writer, templates []*models.Template) {
	templatesTable := util.NewInfoTable(out)
	templatesTable.SetTitle("TEMPLATES")
	templatesTable.AppendHeader(table.Row{"Destination", "Template Path", "Is File", "Description"})
	for _, template := range templates {
		templatesTable.AppendRow(table.Row{template.Destination, template.Path, template.IsFile, template.Description})
	}
	templatesTable.Render()
}

func showPlugins(out io.Writer, plugins []*models.Plugin) {
	pluginsTable := util.NewInfoTable(out)
	pluginsTable.SetTitle("PLUGINS")
	pluginsTable.AppendHeader(table.Row{"Path", "Execution Number", "Description"})
	for _, plugin := range plugins {
		pluginsTable.AppendRow(table.Row{plugin.Path, plugin.ExecNumber, text.WrapSoft(plugin.Description, maxColumnWidth)})
	}
	pluginsTable.Render()
}
