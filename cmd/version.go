package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/version"

	"github.com/spf13/cobra"
)

type versionCommand struct {
	cmd *cobra.Command
}

func newVersionCommand() *versionCommand {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Print the version",
		Aliases:               []string{"v"},
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Proji())
		},
	}
	return &versionCommand{cmd: cmd}
}
