package pkg

import "github.com/spf13/cobra"

// NewCommand returns a new instance of the package command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg"},
		Short:   "Manage packages",
	}

	cmd.AddCommand(
		newImportCommand(),
		newExportCommand(),
		newListCommand(),
		newMimicCommand(),
		newRemoveCommand(),
		newShowCommand(),
	)

	return cmd
}
