package pkg

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"runtime"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
)

func newEditCommand() *cobra.Command {
	var fileType string

	cmd := &cobra.Command{
		Use:                   "edit [OPTIONS] LABEL",
		Short:                 "Edit details about an installed package",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return editPackage(cmd.Context(), args[0], fileType)
		},
	}

	cmd.Flags().StringVarP(&fileType, "type", "t", "toml", "File type in which to apply the editing (toml, json)")

	return cmd
}

func newCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func openFile(ctx context.Context, system, path string) error {
	logger := simplog.FromContext(ctx)

	if path == "" {
		return errors.New("path is empty")
	}

	// TODO: More options for opening files; potentially use config file to set default editor/switch between $EDITOR;
	//       $VISUAL, commands like xdg-open or other application names
	// Set command depending on system
	switch system {
	case "darwin", "freebsd", "linux", "netbsd", "openbsd":
		editor := os.Getenv("EDITOR")
		logger.Debugf("opening file %q with editor %q", path, editor)

		return newCommand(ctx, editor, path).Run()
	case "windows":
		logger.Debugf("opening file %q with start command", path)

		// return exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
		return newCommand(ctx, "start", path).Run()
	default:
		return errors.Newf("unsupported OS %q", system)
	}
}

func editPackage(ctx context.Context, label, fileType string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Load package that should be edited
	logger.Debugf("load package %q", label)
	pkg, err := pama.GetByLabel(ctx, label)
	if err != nil {
		return errors.Wrapf(err, "get package %q", label)
	}

	// Export package to temp dir. Note: an empty destination path will cause the package to be exported to a temp dir
	logger.Debug("exporting package to temp dir")
	path, err := exportPackage(ctx, "", fileType, pkg.ToConfig())
	if err != nil {
		return errors.Wrap(err, "export package")
	}

	// Edit package; open file is OS dependent and will open the file in the default editor.
	logger.Infof("Opening config file of package %q in default editor", label)
	if err := openFile(ctx, runtime.GOOS, path); err != nil {
		return errors.Wrap(err, "open file")
	}

	// Wait for user to finish editing
	logger.Info("Press ENTER to confirm your changes")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// When we are done editing, we can import the edited package
	if err := replacePackage(ctx, label, path); err != nil {
		return errors.Wrap(err, "replace package")
	}

	return nil
}
