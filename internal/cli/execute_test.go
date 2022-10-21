package cli

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/config"
)

func TestMain(m *testing.M) {
	// Set environment to "testing" to help Sentry differentiate between production and testing environments.
	environment = "testing"

	os.Exit(m.Run())
}

func TestExecute(t *testing.T) {
	t.Parallel()

	var cmd *cobra.Command

	// Nil command should be a nop and not panic.
	Execute(cmd)

	// Basic command
	cmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Mute the output
	cmd.SetOut(io.Discard)

	// Execute the command; should not panic
	Execute(cmd)

	// Bind session to cmd context; debug mode
	session := NewSessionWithMode(true)
	ctx := WithSession(context.Background(), session)
	cmd.SetContext(ctx)

	// Replace the command's RunE function with one that returns an error
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return errors.New("test error")
	}

	Execute(cmd)

	// Bind session to cmd context; production mode
	session = NewSessionWithMode(false)
	ctx = WithSession(context.Background(), session)
	cmd.SetContext(ctx)

	// Replace the command's RunE function with one that returns an error
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return errors.New("test error")
	}

	Execute(cmd)

	// Now session with config and sentry disabled
	session.Config = &config.Config{
		Monitoring: config.Monitoring{
			Sentry: config.Sentry{Enabled: false},
		},
	}
	ctx = WithSession(context.Background(), session)
	cmd.SetContext(ctx)

	Execute(cmd)

	// Now session with config and sentry enabled
	session.Config = &config.Config{
		Monitoring: config.Monitoring{
			Sentry: config.Sentry{Enabled: true},
		},
	}
	ctx = WithSession(context.Background(), session)
	cmd.SetContext(ctx)

	Execute(cmd)
}
