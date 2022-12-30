package proji

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/nikoksr/simplog"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/buildinfo"
	"github.com/nikoksr/proji/internal/cli"
	pkg "github.com/nikoksr/proji/internal/cli/proji/package"
	"github.com/nikoksr/proji/internal/cli/proji/server"
	"github.com/nikoksr/proji/internal/cli/proji/version"
	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/internal/manager"
	database "github.com/nikoksr/proji/pkg/database/bolt"
)

var db *database.DB

// cleanup is a helper function that cleans up resources after the CLI has finished. It is meant to be used as by
// the PersistentPostRun function of the root command. We're intentionally not returning errors here, because we want
// to make sure that the cleanup is always executed completely.
func cleanup(ctx context.Context) {
	logger := simplog.FromContext(ctx)

	if db != nil {
		err := db.Close(ctx)
		if err != nil {
			logger.Errorf("failed to close database: %v", err)
		}
	}
}

// rootCommand returns a new instance of the root Command.
func rootCommand() *cobra.Command {
	var debug bool
	var configPath, serverAddress string

	cmd := &cobra.Command{
		Use:           buildinfo.AppName,
		Short:         "A powerful cross-platform CLI project bootstrapper.",
		SilenceErrors: true,
		SilenceUsage:  true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Pick the right logger. Most of the cases, proji will be used as a client, so we'll use the client logger
			// by default and only use the server logger if the user called the server command. The main difference is
			// that the production server logger uses json for structured simplog. We don't want that in the
			// client logger, should be human-readable. Check out the logging package for more info.
			logger := simplog.NewClientLogger(debug)
			if cmd.Name() == "server" {
				logger = simplog.NewServerLogger(debug)
			}

			ctx := simplog.WithLogger(cmd.Context(), logger)

			// Load the app config
			conf, err := config.Load(ctx, configPath, cmd.Flags())
			if err != nil {
				return errors.Wrap(err, "failed to load config")
			}

			// Validate config
			err = conf.Validate()
			if err != nil {
				return errors.Wrap(err, "failed to validate config")
			}

			// Connect to database
			db, err = database.Connect(ctx, conf.Database.DSN)
			if err != nil {
				return errors.Wrap(err, "connect to database")
			}

			// Create package manager
			pama, err := manager.NewPackageManager(ctx, manager.Config{
				Address: serverAddress,
				DB:      db,
				Auth:    &conf.Auth,
				LocalPaths: &manager.LocalPaths{
					Base:      conf.BaseDir(),
					Templates: conf.TemplatesDir(),
					Plugins:   conf.PluginsDir(),
				},
			})
			if err != nil {
				return errors.Wrap(err, "setup package manager")
			}

			// Create project manager
			prma, err := manager.NewProjectManager(ctx, db)
			if err != nil {
				return errors.Wrap(err, "setup project manager")
			}

			// Create a cli.Session and bind it to the command context
			session := cli.NewSessionWithMode(debug).
				WithConfig(conf).
				WithPackageManager(pama).
				WithProjectManager(prma)

			ctx = cli.WithSession(ctx, session)

			// Set the altered context, so that our resources become available to the rest of the application
			cmd.SetContext(ctx)

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cleanup(cmd.Context())
		},
	}

	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	_ = cmd.PersistentFlags().MarkHidden("debug")

	cmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the main config file")
	cmd.PersistentFlags().StringVar(&serverAddress, "remote", "", "Address of a remote proji server to connect to")

	cmd.AddCommand(
		// Client

		// Projects
		projectNewCommand(),
		projectRemoveCommand(),
		projectCleanCommand(),
		projectListCommand(),

		// Packages
		pkg.NewCommand(),

		// Server
		server.NewCommand(),

		// Misc
		version.NewCommand(buildinfo.AppVersion),
	)

	return cmd
}

// Execute is the entrypoint for proji's CLI.
func Execute() {
	cli.Execute(rootCommand())
}
