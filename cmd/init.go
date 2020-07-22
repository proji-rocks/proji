package cmd

import (
	"fmt"
	"runtime"

	"github.com/nikoksr/proji/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type initCommand struct {
	cmd *cobra.Command
}

func newInitCommand() *initCommand {
	cmd := &cobra.Command{
		Use:                   "init",
		Short:                 "Initialize proji",
		Long:                  initHelp(),
		Hidden:                true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Deploy()
			if err != nil {
				return errors.Wrap(err, "could not set up config folder")
			}
			return nil
		},
	}
	return &initCommand{cmd: cmd}
}

func initHelp() string {
	// OS specific conf path.
	confPath := config.GetBaseConfigPath()

	// OS specific help command.
	helpMsg := "It is possible to set up the config folder manually.\n\n"

	switch runtime.GOOS {
	case "darwin", "linux":
		helpMsg += fmt.Sprintf(
			" mkdir -p %s/db %s/examples %s/plugins %s/templates\n\n",
			confPath,
			confPath,
			confPath,
			confPath,
		)
		helpMsg += fmt.Sprintf(
			" curl -o %s/config.toml https://raw.githubusercontent.com/nikoksr/proji/master/assets/examples/example-config.toml\n\n",
			confPath,
		)
		helpMsg += fmt.Sprintf(
			" curl -o %s/examples/proji-package.toml https://raw.githubusercontent.com/nikoksr/proji/master/assets/examples/example-package-export.toml\n\n",
			confPath,
		)
	case "windows":
		helpMsg += fmt.Sprintf(
			" md %s\\db %s\\examples %s\\plugins %s\\templates\n\n",
			confPath,
			confPath,
			confPath,
			confPath,
		)
		helpMsg += fmt.Sprintf(
			" Download https://github.com/nikoksr/proji/blob/master/assets/examples/example-config.toml to %s\\config.toml\n",
			confPath,
		)
		helpMsg += fmt.Sprintf(
			" Download https://github.com/nikoksr/proji/blob/master/assets/examples/example-package-export.toml to %s\\examples\\proji-package.toml\n\n",
			confPath,
		)
	default:
		helpMsg = "Your platform is not supported, so no help is available at the moment.\n\n"
	}
	helpMsg += "\nFor more help visit: https://github.com/nikoksr/proji\n\n"
	return helpMsg
}
