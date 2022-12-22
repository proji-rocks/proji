package pkg

import (
	"context"
	"strings"

	"github.com/nikoksr/proji/pkg/pointer"

	"github.com/nikoksr/proji/pkg/api/v1/domain"

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
	var example bool

	cmd := &cobra.Command{
		Use:                   "export [OPTIONS] LABEL [LABEL...]",
		Short:                 "Export packages as config files",
		Aliases:               []string{"out"},
		DisableFlagsInUseLine: true,

		Example: `  proji package export py
  proji package out py js
  proji package out -d ./my_packages cpp go
  proji package out -t json py3`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// If example flag is set, export only the example package. In this case, no labels are required.
			if example {
				return exportExample(cmd.Context(), destination, fileType)
			}

			// If no labels are provided, return an error
			if len(args) == 0 {
				return errors.New("no labels provided")
			}

			return exportPackages(cmd.Context(), destination, fileType, args...)
		},
	}

	cmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination folder")
	cmd.Flags().StringVarP(&fileType, "type", "t", "toml", "File type to export to (toml, json)")
	cmd.Flags().BoolVarP(&example, "example", "e", false, "Export example config file")

	return cmd
}

func exportPackage(ctx context.Context, destination, fileType string, pkg *domain.PackageConfig) (path string, err error) {
	switch strings.ToLower(fileType) {
	case portability.FileTypeTOML:
		path, err = exporting.ToTOML(ctx, pkg, destination)
	case portability.FileTypeJSON:
		path, err = exporting.ToJSON(ctx, pkg, destination)
	default:
		err = portability.ErrUnsupportedConfigFileType
	}

	return path, err
}

func exportExample(ctx context.Context, destination, fileType string) error {
	logger := simplog.FromContext(ctx)

	examplePkg := domain.PackageConfig{
		Label:       "xxx",
		Name:        "Example",
		Description: pointer.To("This is an example package"),
		DirTree: &domain.DirTreeConfig{
			Entries: []*domain.DirEntryConfig{
				{IsDir: true, Path: "docs"},
				{IsDir: false, Path: "docs/docs.md", Template: &domain.TemplateConfig{
					Path: "github/nikoksr/docs.md",
				}},
				{IsDir: true, Path: "src"},
				{IsDir: false, Path: "src/main.go", Template: &domain.TemplateConfig{
					Path: "github/nikoksr/main.go",
				}},
				{IsDir: false, Path: "1_You"},
				{IsDir: false, Path: "2_Are"},
				{IsDir: false, Path: "3_Awesome"},
			},
		},
		Plugins: &domain.PluginSchedulerConfig{
			Pre:  []*domain.PluginConfig{{Path: "github/nikoksr/go-init.lua"}},
			Post: []*domain.PluginConfig{{Path: "github/nikoksr/git-init.lua"}},
		},
	}

	// Export package
	path, err := exportPackage(ctx, destination, fileType, &examplePkg)
	if err != nil {
		return errors.Wrap(err, "export package")
	}

	logger.Infof("Exported example package to %q", path)

	return nil
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
		path, err := exportPackage(ctx, destination, fileType, pkg.ToConfig())
		if err != nil {
			logger.Errorf("Failed to export package %q: %v", label, err)
		} else {
			logger.Infof("Exported package %q to %q", label, path)
		}
	}

	return nil
}
