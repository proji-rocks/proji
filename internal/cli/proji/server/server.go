package server

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/internal/server"
)

// NewCommand returns a new instance of the server command.
func NewCommand() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run as server",

		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(cmd.Context(), address)
		},
	}

	cmd.Flags().StringVarP(&address, "address", "a", ":8080", "Address to listen on")

	cmd.AddCommand(newPingCommand())

	return cmd
}

func serve(ctx context.Context, address string) error {
	// Load session to get DSN from config.
	session := cli.SessionFromContext(ctx)
	if session == nil {
		return errors.New("cli session is nil")
	}
	if session.Config == nil {
		return errors.New("cli session config is nil")
	}

	// Create the server.
	srvr, err := server.New(ctx, &server.Managers{
		Package: session.PackageManager,
	})
	if err != nil {
		return errors.Wrap(err, "create server")
	}

	// And run it. It will block until the context is canceled. It will shut down gracefully.
	return srvr.Run(ctx, address)
}
