package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand returns a new instance of the version command.
func NewCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the app's version",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
