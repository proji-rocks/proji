//nolint:gochecknoglobals,gochecknoinits
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

// Session represents central resources and information the app uses.
type Session struct {
	Config          *config.Config
	StorageService  storage.Service
	FallbackVersion string
	Version         string
	NoColors        bool
}

var session *Session
var terminalWidth, maxColumnWidth int
var disableColors bool

const (
	defaultMaxColumnWidth = 50
)

var rootCmd = &cobra.Command{
	Use:           "proji",
	Short:         "A powerful cross-platform CLI project templating tool.",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		messages.EnableColors(disableColors)

		// Leave one empty line above by default
		fmt.Println()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		messages.Error("", err)
	}
}

func init() {
	if session == nil {
		session = &Session{
			Config:          nil,
			StorageService:  nil,
			FallbackVersion: "0.19.2",
			Version:         "0.20.0",
		}
	}

	// Initialize command flags
	initFlags()

	// Skip initialization if no args were given
	if len(os.Args) < 2 {
		return
	}

	// Evaluate initialization behaviour
	var initFunctions []func()
	switch os.Args[1] {
	case "version", "help":
		// Don't init config or storage on version or help. It's just not necessary.
		return
	case "init":
		// Setup the config because init needs a bare bone config to deploy the base config folder.
		initFunctions = append(initFunctions, setupConfig)
	default:
		// On default load the main config and initialize the storage service
		initFunctions = []func(){
			initConfig,
			initStorageService,
		}
	}
	cobra.OnInitialize(initFunctions...)
}

func initFlags() {
	rootCmd.PersistentFlags().BoolVar(&disableColors, "no-colors", false, "disable text colors")
}

func setupConfig() {
	err := config.Setup()
	if err != nil {
		messages.Error("failed to setup config", err)
		os.Exit(1)
	}
}

func initConfig() {
	if session == nil {
		messages.Error("couldn't initialize config", fmt.Errorf("session not found"))
	}

	// Run config setup
	setupConfig()

	// Create the config
	session.Config = config.New(config.GetBaseConfigPath())

	// Load the config
	err := session.Config.Load()
	if err != nil {
		messages.Error("loading config failed", err)
		os.Exit(1)
	}
}

func initStorageService() {
	var err error
	session.StorageService, err = storage.NewService(
		session.Config.DatabaseConnection.Driver,
		session.Config.DatabaseConnection.DSN,
	)
	if err != nil {
		messages.Error(
			"could not connect to %s database with dsn %s, %s",
			err,
			session.Config.DatabaseConnection.Driver,
			session.Config.DatabaseConnection.DSN,
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

func setMaxColumnWidth() {
	//Load terminal width and set max column width for dynamic rendering
	var err error
	terminalWidth, err = getTerminalWidth()
	if err != nil {
		messages.Warning("couldn't get terminal width. Falling back to default value, %s", err.Error())
		maxColumnWidth = defaultMaxColumnWidth
	} else {
		maxColumnWidth = terminalWidth / 2
	}
}
