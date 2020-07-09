package cmd

import (
	"os"

	"github.com/nikoksr/proji/static"

	"github.com/spf13/cobra"
)

type completionCommand struct {
	cmd *cobra.Command
}

func newCompletionCommand() *completionCommand {
	var cmd = &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Load shell completions",
		Long:                  static.CompletionHelpMessage,
		Hidden:                true,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
			return err
		},
	}
	return &completionCommand{cmd: cmd}
}
