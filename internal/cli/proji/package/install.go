package pkg

import (
	"context"
	"regexp"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/packages/portability/importing"
	"github.com/nikoksr/proji/pkg/remote"
)

func newImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "install [OPTIONS] PATH [PATH...]",
		Short:                 "Install packages from local or remote config files",
		Aliases:               []string{"i"},
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,

		Example: `  proji package install https://github.com/nikoksr/my_repo/blob/main/my_package.json
  proji package install gh://nikoksr/my_repo/blob/main/my_package.json
  proji package in gh://nikoksr/my_repo/blob/main/my_package.json
  proji package in /home/my_user/my_package.json`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return installPackages(cmd.Context(), args...)
		},
	}

	return cmd
}

type pathType int

const (
	pathTypeURL pathType = iota
	pathTypeLocal
)

func getPathType(path string) pathType {
	parsedPath, err := remote.ParseRepoURL(path)

	if err == nil && parsedPath != nil && parsedPath.Scheme != "" && parsedPath.Host != "" {
		return pathTypeURL
	}

	return pathTypeLocal
}

func importPackage(ctx context.Context, path string, isFile bool, exclude *regexp.Regexp) (*domain.PackageAdd, error) {
	var _package *domain.PackageAdd
	var err error

	switch getPathType(path) {
	case pathTypeURL:
		if isFile {
			_package, err = importing.RemotePackage(ctx, path)
		} else {
			_package, err = importing.RepositoryAsPackage(ctx, path, exclude)
		}
	case pathTypeLocal:
		if isFile {
			_package, err = importing.LocalPackage(ctx, path)
		} else {
			_package, err = importing.LocalFolderAsPackage(ctx, path, exclude)
		}
	default:
		return nil, errors.Wrapf(err, "invalid path %q", path)
	}

	return _package, err
}

func installPackages(ctx context.Context, paths ...string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Install packages
	logger.Debugf("installing %d packages", len(paths))
	for _, path := range paths {
		if path == "" {
			logger.Debugf("skipping import due to empty path")
			continue
		}

		logger.Debugf("importing package from %q", path)
		_package, err := importPackage(ctx, path, true, nil)
		if err != nil {
			return errors.Wrapf(err, "import package from %q", path)
		}

		logger.Debugf("adding package %q", _package.Label)
		if err = pama.Store(ctx, _package); err != nil {
			return errors.Wrapf(err, "store %q, imported from %q", _package.Name, path)
		}

		logger.Infof("Successfully installed package %q as %q", _package.Name, _package.Label)
	}

	return nil
}
