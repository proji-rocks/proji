//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"log"
	"os"

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
}

var session *Session
var terminalWidth, maxColumnWidth int

const (
	defaultMaxColumnWidth = 50
)

var rootCmd = &cobra.Command{
	Use:   "proji",
	Short: "A powerful cross-platform CLI project templating tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if session == nil {
			log.Fatalln("session is not defined")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = rootCmd.Execute()
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
		// Setup the config because init needs a barebone config to deploy the base config folder.
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

func setupConfig() {
	err := config.Setup()
	if err != nil {
		log.Fatalf("failed to setup config, %s", err.Error())
	}
}

func initConfig() {
	if session == nil {
		log.Fatalf("couldn't set config, environment struct is nil")
	}

	// Run config setup
	setupConfig()

	// Create the config
	session.Config = config.New(config.GetBaseConfigPath())

	// Load the config
	err := session.Config.Load()
	if err != nil {
		log.Fatalf("loading config failed, %s", err.Error())
	}
}

func initStorageService() {
	var err error
	session.StorageService, err = storage.NewService(
		session.Config.DatabaseConnection.Driver,
		session.Config.DatabaseConnection.DSN,
	)
	if err != nil {
		log.Fatalf(
			"Error: could not connect to %s database with dsn %s, %s\n",
			session.Config.DatabaseConnection.Driver,
			session.Config.DatabaseConnection.DSN,
			err.Error(),
		)
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
		fmt.Printf("Warning: Couldn't get terminal width, %s", err.Error())
		maxColumnWidth = defaultMaxColumnWidth
	} else {
		maxColumnWidth = terminalWidth / 2
	}
}
