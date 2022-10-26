package server

import (
	"context"
	"strings"
	"time"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/pkg/sdk/health"
)

func newPingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping ADDRESS",
		Short: "Ping a remote package server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[0]

			return ping(cmd.Context(), address)
		},
	}

	return cmd
}

func pingRemoteServer(ctx context.Context, address string) (time.Duration, error) {
	client, err := health.NewClient(address)
	if err != nil {
		return -1, errors.Wrap(err, "create client")
	}

	// isHealthy is false if err is not nil, so we don't need to check err here explicitly.
	start := time.Now()
	isHealthy, err := client.IsHealthy(ctx)
	latency := time.Since(start)
	if isHealthy {
		return latency, nil
	}

	// Check for common errors.
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "no route to host") {
		return -1, errors.New("not reachable")
	} else if strings.Contains(err.Error(), "connection timed out") {
		return -1, errors.New("connection timed out")
	}

	return -1, err
}

func ping(ctx context.Context, address string) error {
	logger := simplog.FromContext(ctx)

	// Sanitize and check address.
	address = strings.ToLower(strings.TrimSpace(address))
	if address == "" {
		return errors.New("Server address is empty")
	}
	/*	if address == "localhost" || address == "127.0.0.1" {
			return errors.New("Server address is localhost")
		}
	*/
	// Ping the server.
	logger.Infof("Pinging package server %q", address)
	latency, err := pingRemoteServer(ctx, address)
	if err != nil {
		logger.Warnf("Failed to ping package server %q: %v", address, err)
	} else {
		logger.Infof("Package server %q responded in %v", address, latency)
	}

	return nil
}
