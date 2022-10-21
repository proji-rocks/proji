package plugins

import (
	"context"
	"os"
	"os/exec"

	"github.com/nikoksr/proji/pkg/logging"
)

// TODO: This needs Windows support - sigh.
func run(ctx context.Context, path string) error {
	logger := logging.FromContext(ctx)

	cmd := exec.Command("lua", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	logger.Debugf("executing lua script %s", path)

	return cmd.Run()
}

// Run runs the lua script at path. It is equivalent to: `lua <path>`.
func Run(ctx context.Context, path string) error {
	return run(ctx, path)
}
