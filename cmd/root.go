package cmd

import (
	"os"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/internal/database"
	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/internal/statuswriter"
	"github.com/nikoksr/proji/pkg/domain"
	packageservice "github.com/nikoksr/proji/pkg/package/service"
	packagestore "github.com/nikoksr/proji/pkg/package/store"
	projectservice "github.com/nikoksr/proji/pkg/project/service"
	projectstore "github.com/nikoksr/proji/pkg/project/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

//nolint:gochecknoglobals
var session *sessionState

// sessionState represents central resources and information the app uses.
type sessionState struct {
	config              *config.Config
	packageService      domain.PackageService
	projectService      domain.ProjectService
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
		message.Errorf("", err)
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if disableColors {
				message.DisableColors()
			}

			// Leave one empty line above by default
			// fmt.Println()

			// Prepare proji
			return prepare()
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

func prepare() error {
	if session == nil {
		session = &sessionState{
			fallbackVersion:     "0.19.2",
			version:             "0.20.0",
			noColors:            false,
			maxTableColumnWidth: getMaxColumnWidth(),
		}
	}

	// Skip preparation if no args were given
	if len(os.Args) < 2 {
		return nil
	}

	// Evaluate preparation behaviour
	var err error
	switch os.Args[1] {
	case "version", "help":
		// Do nothing. Don't init config or storage on version or help. It's just not necessary.
	case "init":
		// Prepare the config because init needs a bare bone config to deploy the base config folder.
		err = config.Prepare()
	default:
		// On default load the main config and initialize the storage service
		err = loadConfig()
		if err != nil {
			return err
		}
		err = initServices()
	}
	return err
}

func loadConfig() error {
	if session == nil {
		message.Errorf("could not initialize config", fmt.Errorf("session not found"))
	}

	// Prepare the config
	err := config.Prepare()
	if err != nil {
		return errors.Wrap(err, "prepare config")
	}

	// Create the config
	session.config = config.New(config.GetBaseConfigPath())

	// Load the config
	err = session.config.LoadValues()
	if err != nil {
		return errors.Wrap(err, "load config values")
	}
	return nil
}

func initServices() error {
	// Connect to database
	db, err := database.New(session.config.DatabaseConnection.Driver, session.config.DatabaseConnection.DSN)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}

	// Run database migration
	err = db.Migrate()
	if err != nil {
		return errors.Wrap(err, "migrate database")
	}

	// Create the services
	createServices(db)
	return nil
}

func createServices(db *database.Database) {
	// Package Service
	packageStore := packagestore.New(db.Connection)
	session.packageService = packageservice.New(session.config.Auth, packageStore)

	// Project Service
	projectStore := projectstore.New(db.Connection)
	session.projectService = projectservice.New(projectStore)
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
		message.Warningf("couldn't get terminal width. Falling back to default value, %s", err.Error())
		return 50
	}
	return terminalWidth / 2
}
