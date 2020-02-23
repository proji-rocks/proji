package cmd

import (
	"log"

	"github.com/nikoksr/proji/pkg/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize user-specific config folder",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		projiEnv.UserConfigPath, err = config.InitConfig(projiEnv.UserConfigPath)
		if err != nil {
			// TODO: Improve error message. Manual config setup is possible.
			log.Fatalf("Error: could not set up config folder. %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
