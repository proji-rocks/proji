package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/internal/message"
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
	err := config.Prepare()
	if err != nil {
		log.Println(message.Serrorf(err, "failed to prepare config"))
	}
	configDirectoryPath := config.GetBaseConfigPath()
	configPath := filepath.Join(configDirectoryPath, "config.toml")
	dbPath := filepath.Join(configDirectoryPath, "db")
	pluginsPath := filepath.Join(configDirectoryPath, "plugins")
	templatesPath := filepath.Join(configDirectoryPath, "templates")
	helpMsg := "In the case that proji's initialization fails you can create its central config folder manually.\n\n"

	switch runtime.GOOS {
	case "darwin", "linux":
		helpMsg += fmt.Sprintf(" • mkdir -p %s %s %s\n",
			dbPath,
			pluginsPath,
			templatesPath,
		)
		helpMsg += fmt.Sprintf(
			" • curl https://github.com/nikoksr/proji/examples/main-config.toml -o %s\n",
			configPath,
		)
	case "windows":
		helpMsg += fmt.Sprintf(
			" • md %s %s %s\n",
			dbPath,
			pluginsPath,
			templatesPath,
		)
		helpMsg += " • Download the main config from: https://github.com/nikoksr/proji/examples/main-config.toml\n"
		helpMsg += fmt.Sprintf(" • Move it to: %s\n", configPath)
	default:
		return "Your operating system is not supported, sorry!\n"
	}
	helpMsg += "\nFor more help visit: https://github.com/nikoksr/proji\n\n"
	return helpMsg
}
