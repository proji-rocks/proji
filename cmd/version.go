package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type versionCommand struct {
	cmd *cobra.Command
}

func newVersionCommand() *versionCommand {
	var cmd = &cobra.Command{
		Use:                   "version",
		Short:                 "Print the version",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v" + activeSession.version)
		},
	}
	return &versionCommand{cmd: cmd}
}
