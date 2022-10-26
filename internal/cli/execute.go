package cli

import (
	"context"
	"time"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/buildinfo"
)

const sentryDSN = "https://c02c90cfaec14f0c86247caee1c7de7b@o408463.ingest.sentry.io/6340441"

// environment is used to determine the environment in which the application is running. We pass it to Sentry to help
// us determine the environment in which the error occurred. This defaults to "production" but can be overridden by
// TestMain in execute_test.go, to allow for differentiation between production and testing environments.
var environment = "production"

// Execute is a function that abstracts a lot of the boilerplate code that is needed to execute a command. It is meant
// to be used in the main function of a command. It will handle errors, logging, and sentry reporting. If cmd is nil,
// it's a nop.
func Execute(cmd *cobra.Command) {
	if cmd == nil {
		return
	}

	// Actual execution of the root command and entrypoint for proji. Using ExecuteC() instead of Execute() here because
	// we want to access the session that was attached to the root command's context.
	cmd, err := cmd.ExecuteC()

	// Cleanly handle possible errors.
	handleExecutionError(cmd.Context(), err)
}

func handleExecutionError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	// Logger and session should've been set in the root command's PersistentPreRunE function.
	logger := simplog.FromContext(ctx)

	// We need the session to check if we're in debug mode and if reporting to Sentry is enabled.
	session := SessionFromContext(ctx)

	// Log the error.
	if session.Debug {
		logger.Errorf("%+v", err) // This prints the stacktrace created by cockroachdb/errors.
	} else {
		logger.Error(err)
	}

	// Config holds the Sentry configuration. So if it's not set, we don't need to report to Sentry.
	if session.Config == nil {
		logger.Debugf("no config found in session")
		return
	}

	// Only report errors to Sentry if enabled in settings, and we're not in debug mode. Dirty builds are interpreted as
	// development builds, so we don't report them to Sentry.
	if !session.Config.Monitoring.Sentry.Enabled || session.Debug || buildinfo.BuildDirty {
		logger.Debugf(
			"Not reporting to Sentry - Enabled: %t, Debug: %t, Dirty: %t",
			session.Config.Monitoring.Sentry.Enabled, session.Debug, buildinfo.BuildDirty,
		)

		return
	}

	// If we are here, we are allowed to monitor, initialize Sentry.
	initErr := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Release:          buildinfo.AppVersion,
		Environment:      environment,
		AttachStacktrace: true,
		Transport:        sentry.NewHTTPTransport(),
	})
	if initErr != nil {
		logger.Errorf("Failed to initialize Sentry client: %v", initErr)
		return
	}

	// Sentry initialization successful, report the error.
	eventID := errors.ReportError(err)
	if eventID == "" {
		logger.Warn("Failed to report error to Sentry")
		return
	}

	logger.Infof("Created Sentry error event with ID: %q", eventID)
	if !sentry.Flush(2 * time.Second) {
		logger.Warn("Flushing Sentry client timed out, some events may have been lost")
	} else {
		logger.Info("Successfully sent error to Sentry. Thanks for the feedback!")
	}
}
