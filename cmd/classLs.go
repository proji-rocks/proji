//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"io"
	"os"

	"github.com/nikoksr/proji/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command.
var classLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listClasses(os.Stdout)
	},
}

func init() {
	classCmd.AddCommand(classLsCmd)
}

func listClasses(out io.Writer) error {
	classes, err := projiEnv.StorageService.LoadClasses()
	if err != nil {
		return err
	}

	classesTable := util.NewInfoTable(out)
	classesTable.AppendHeader(table.Row{"Name", "Label"})

	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		classesTable.AppendRow(table.Row{class.Name, class.Label})
	}
	classesTable.Render()
	return nil
}
