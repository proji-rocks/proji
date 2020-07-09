package cmd

import (
	"fmt"
	"os"

	"github.com/nikoksr/proji/messages"

	"github.com/nikoksr/proji/config"
	"github.com/nikoksr/proji/storage"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

//nolint:gochecknoglobals
var activeSession *session

// session represents central resources and information the app uses.
type session struct {
	config              *config.Config
	storageService      storage.Service
	fallbackVersion     string
	version             string
	noColors            bool
	maxTableColumnWidth int
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCommand().cmd.Execute()
	if err != nil {
		messages.Errorf("", err)
	}
}

type rootCommand struct {
	cmd *cobra.Command
}

func newRootCommand() *rootCommand {
	var disableColors bool

	var cmd = &cobra.Command{
		Use:           "proji",
		Short:         "A powerful cross-platform CLI project templating tool.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if disableColors {
				messages.DisableColors()
			}

			// Leave one empty line above by default
			fmt.Println()

			// Prepare proji
			prepare()
		},
	}

	cmd.PersistentFlags().BoolVar(&disableColors, "no-colors", false, "disable text colors")
	cmd.AddCommand(
		newCompletionCommand().cmd,
		newInitCommand().cmd,
		newPackageCommand().cmd,
		newProjectAddCommand().cmd,
		newProjectCleanCommand().cmd,
		newProjectCreateCommand().cmd,
		newProjectListCommand().cmd,
		newProjectRemoveCommand().cmd,
		newProjectSetCommand().cmd,
		newVersionCommand().cmd,
	)
	return &rootCommand{cmd: cmd}
}

func prepare() {
	if activeSession == nil {
		activeSession = &session{
			config:              nil,
			storageService:      nil,
			fallbackVersion:     "0.19.2",
			version:             "0.20.0",
			noColors:            false,
			maxTableColumnWidth: getMaxColumnWidth(),
		}
	}

	// Skip preparation if no args were given
	if len(os.Args) < 2 {
		return
	}

	// Evaluate preparation behaviour
	switch os.Args[1] {
	case "version", "help":
		// Don't init config or storage on version or help. It's just not necessary.
		return
	case "init":
		// Setup the config because init needs a bare bone config to deploy the base config folder.
		setupConfig()
	default:
		// On default load the main config and initialize the storage service
		loadConfig()
		initStorageService()
	}
}

func setupConfig() {
	err := config.Setup()
	if err != nil {
		messages.Errorf("failed to setup config", err)
		os.Exit(1)
	}
}

func loadConfig() {
	if activeSession == nil {
		messages.Errorf("couldn't initialize config", fmt.Errorf("session not found"))
	}

	// Run config setup
	setupConfig()

	// Create the config
	activeSession.config = config.New(config.GetBaseConfigPath())

	// Load the config
	err := activeSession.config.Load()
	if err != nil {
		messages.Errorf("loading config failed", err)
		os.Exit(1)
	}
}

func initStorageService() {
	var err error
	activeSession.storageService, err = storage.NewService(
		activeSession.config.DatabaseConnection.Driver,
		activeSession.config.DatabaseConnection.DSN,
	)
	if err != nil {
		messages.Errorf(
			"could not connect to %s database with dsn %s, %s",
			err,
			activeSession.config.DatabaseConnection.Driver,
			activeSession.config.DatabaseConnection.DSN,
		)
		os.Exit(1)
	}
}

func getTerminalWidth() (int, error) {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, err
	}
	return w, nil
}

func getMaxColumnWidth() int {
	//Load terminal width and set max column width for dynamic rendering
	terminalWidth, err := getTerminalWidth()
	if err != nil {
		messages.Warningf("couldn't get terminal width. Falling back to default value, %s", err.Error())
		return 50
	}
	return terminalWidth / 2
}
