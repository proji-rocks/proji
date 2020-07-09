package cmd

import (
	"github.com/spf13/cobra"
)

type projectSetCommand struct {
	cmd *cobra.Command
}

func newProjectSetCommand() *projectSetCommand {
	var cmd = &cobra.Command{
		Use:   "set",
		Short: "Set project information",
	}

	cmd.AddCommand(newProjectSetPathCommand().cmd)

	return &projectSetCommand{cmd: cmd}
}
