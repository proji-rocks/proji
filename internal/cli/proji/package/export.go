package pkg

import (
	"context"
	"strings"

	"github.com/nikoksr/proji/pkg/packages/portability"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/pkg/packages/portability/exporting"
)

func newExportCommand() *cobra.Command {
	var destination string
	var fileType string

	cmd := &cobra.Command{
		Use:                   "export [OPTIONS] LABEL [LABEL...]",
		Short:                 "Export packages as config files",
		Aliases:               []string{"out"},
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,

		Example: `  proji package export py
  proji package out py js
  proji package out -d ./my_packages cpp go
  proji package out -t json py3`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return exportPackages(cmd.Context(), destination, fileType, args...)
		},
	}

	cmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination folder")
	cmd.Flags().StringVarP(&fileType, "type", "t", "toml", "File type to export to (toml, json)")

	return cmd
}

func exportPackages(ctx context.Context, destination, fileType string, labels ...string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Load packages
	logger.Debugf("exporting %d packages", len(labels))
	for _, label := range labels {
		logger.Debugf("loading package %q", label)
		pkg, err := pama.GetByLabel(ctx, label)
		if err != nil {
			return errors.Wrapf(err, "get package %q", label)
		}

		// Export package
		logger.Debugf("exporting package %q as %q to %q", label, fileType, destination)

		var path string
		switch strings.ToLower(fileType) {
		case portability.FileTypeTOML:
			path, err = exporting.ToTOML(ctx, pkg.AsExportable(), destination)
		case portability.FileTypeJSON:
			path, err = exporting.ToJSON(ctx, pkg.AsExportable(), destination)
		default:
			err = portability.ErrUnsupportedConfigFileType
		}

		if err != nil {
			logger.Errorf("Failed to export package %q: %v", label, err)
		} else {
			logger.Infof("Exported package %q to %q", label, path)
		}
	}

	return nil
}
