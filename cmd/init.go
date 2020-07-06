//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"log"
	"runtime"

	"github.com/nikoksr/proji/pkg/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize user-specific config folder",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitConfig(
			projiEnv.ConfigFolderPath,
			projiEnv.Version,
			projiEnv.FallbackVersion,
			false,
		)
		if err != nil {
			log.Fatalf(
				"Error: could not set up config folder. %v\n\n"+initHelp(),
				err,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initHelp() string {
	// OS specific conf path.
	confPath, err := config.GetBaseConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	// OS specific help command.
	helpMsg := "It is possible to set up the config folder manually.\n\n"

	switch runtime.GOOS {
	case "darwin", "linux":
		helpMsg += fmt.Sprintf(
			" mkdir -p %s/db %s/examples %s/scripts %s/templates\n\n",
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
			" curl -o %s/examples/proji-class.toml https://raw.githubusercontent.com/nikoksr/proji/master/assets/examples/example-class-export.toml\n\n",
			confPath,
		)
	case "windows":
		helpMsg += fmt.Sprintf(
			" md %s\\db %s\\examples %s\\scripts %s\\templates\n\n",
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
			" Download https://github.com/nikoksr/proji/blob/master/assets/examples/example-class-export.toml to %s\\examples\\proji-class.toml\n\n",
			confPath,
		)
	default:
		helpMsg = "Your platform is not supported, so no help is available at the moment.\n\n"
	}
	helpMsg += "\nFor more help visit: https://github.com/nikoksr/proji\n\n"
	return helpMsg
}
