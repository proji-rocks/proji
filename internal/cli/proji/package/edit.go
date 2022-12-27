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

func openFile(ctx context.Context, system, editor, path string) error {
	logger := simplog.FromContext(ctx)

	if path == "" {
		return errors.New("path is empty")
	}

	if editor == "" {
		switch system {
		case "darwin", "freebsd", "linux", "netbsd", "openbsd":
			editor = os.Getenv("EDITOR")
		case "windows":
			editor = "notepad.exe"
		default:
			return errors.Newf("unsupported OS %q", system)
		}
	}

	logger.Debugf("opening file %q in %s", path, editor)

	return newCommand(ctx, editor, path).Run()
}

func editPackage(ctx context.Context, label, fileType string) error {
	logger := simplog.FromContext(ctx)

	// Get the session
	logger.Debug("getting package manager from cli session")
	session := cli.SessionFromContext(ctx)

	// Get the package manager
	pama := session.PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Get value for the text editor defined in the config file
	textEditor := session.Config.System.TextEditor

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
	logger.Infof("Opening config of %q", label)
	if err := openFile(ctx, runtime.GOOS, textEditor, path); err != nil {
		return errors.Wrap(err, "open file")
	}

	// Wait for user to finish editing
	logger.Info("Press ENTER to confirm your changes when you are done")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

	// When we are done editing, we can import the edited package
	if err := replacePackage(ctx, label, path); err != nil {
		return errors.Wrap(err, "replace package")
	}

	return nil
}
