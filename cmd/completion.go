//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(proji completion bash)

# To load completions for each session, execute once:
Linux:
  $ proji completion bash > /etc/bash_completion.d/proji
MacOS:
  $ proji completion bash > /usr/local/etc/bash_completion.d/proji

Zsh:

$ source <(proji completion zsh)

# To load completions for each session, execute once:
$ proji completion zsh > "${fpath[1]}/_proji"

Fish:

$ proji completion fish | source

# To load completions for each session, execute once:
$ proji completion fish > ~/.config/fish/completions/proji.fish
`,
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

func init() {
	rootCmd.AddCommand(completionCmd)
}
