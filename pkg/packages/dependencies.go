package packages

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/nikoksr/simplog"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/nikoksr/proji/pkg/remote/platform"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)

	return errors.Is(err, os.ErrNotExist)
}

func (m *localManager) downloadDependency(ctx context.Context, upstreamURL *string, basePath string) error {
	if upstreamURL == nil || *upstreamURL == "" {
		return errors.New("upstreamURL is empty")
	}
	if basePath == "" {
		return errors.New("Base path is empty")
	}

	// Make sure the Base path is cross-platform compatible.
	basePath = filepath.FromSlash(basePath)

	// Extract information about the dependency from the upstream URL
	upstream, err := url.Parse(*upstreamURL)
	if err != nil {
		return errors.Wrap(err, "parse upstreamURL")
	}

	repoInfo, err := remote.ExtractRepoInfoFromURL(ctx, upstream)
	if err != nil {
		return errors.Wrap(err, "extract repo information")
	}

	// Create the platform based on the upstream host.
	_platform, err := platform.NewWithAuth(ctx, upstream.Hostname(), m.auth)
	if err != nil {
		return errors.Wrap(err, "identify platform")
	}

	// Set together the full destination path.
	//
	// We concat the storage path, with the platform name, the owner's name and the dependency's name. The plugin's
	// name is the last part of the upstream URL.
	//
	// For example, if the upstream URL is: https://github.com/nikoksr/proji/blob/main/plugins/plugin.lua
	//
	// The destination path will be: <base_path>/github/nikoksr/plugin.lua
	owner := repoInfo.Owner
	name := path.Base(upstream.Path)
	destination := filepath.Join(basePath, _platform.String(), owner, name)

	// Only download the dependency if it doesn't exist yet.
	if pathExists(destination) {
		return os.ErrExist
	}

	// Download the plugin to the destination.
	return _platform.DownloadFile(ctx, repoInfo, *upstreamURL, destination)
}

func (m *localManager) downloadPlugin(ctx context.Context, plugin *domain.Plugin) error {
	logger := simplog.FromContext(ctx)
	logger.Debugf("downloading plugin %v", plugin.UpstreamURL)

	return m.downloadDependency(ctx, plugin.UpstreamURL, m.paths.Plugins)
}

func (m *localManager) downloadTemplate(ctx context.Context, template *domain.Template) error {
	logger := simplog.FromContext(ctx)
	logger.Debugf("downloading template %v", template.UpstreamURL)

	return m.downloadDependency(ctx, template.UpstreamURL, m.paths.Templates)
}

func (m *localManager) downloadDependencies(ctx context.Context, pkg *domain.PackageAdd) (err error) {
	if pkg == nil {
		return errors.New("package is nil")
	}

	// Download Templates
	for _, entry := range pkg.DirTree.Entries {
		if entry == nil || entry.Template == nil || entry.Template.UpstreamURL == nil {
			continue
		}

		if err = m.downloadTemplate(ctx, entry.Template); err != nil {
			return errors.Wrap(err, "download template")
		}
	}

	// Download Plugins
	if pkg.Plugins == nil {
		return nil
	}

	for _, plugin := range append(pkg.Plugins.Pre, pkg.Plugins.Post...) {
		if plugin == nil || plugin.UpstreamURL == nil {
			continue
		}
		if err = m.downloadPlugin(ctx, plugin); err != nil {
			return errors.Wrap(err, "download plugin")
		}
	}

	return nil
}
