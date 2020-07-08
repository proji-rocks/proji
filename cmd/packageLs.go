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
var packageLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listPackages(os.Stdout)
	},
}

func init() {
	packageCmd.AddCommand(packageLsCmd)
}

func listPackages(out io.Writer) error {
	packages, err := session.StorageService.LoadPackages()
	if err != nil {
		return err
	}

	packagesTable := util.NewInfoTable(out)
	packagesTable.AppendHeader(table.Row{"Name", "Label"})

	for _, pkg := range packages {
		if pkg.IsDefault {
			continue
		}
		packagesTable.AppendRow(table.Row{pkg.Name, pkg.Label})
	}
	packagesTable.Render()
	return nil
}
