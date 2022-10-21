package importing

import (
	"bytes"
	"context"
	"regexp"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/nikoksr/proji/pkg/remote/platform"
)

//
// Utility functions
//

// regexToPathSkipper returns a function that can be used to skip a path if it matches the given regex.
func regexToPathSkipper(regex *regexp.Regexp) remote.PathSkipperFn {
	if regex == nil || regex.String() == "" {
		return remote.DefaultPathSkipper
	}

	return func(path string) bool {
		return regex.MatchString(path)
	}
}

//
// Remote Importing
//

// RemotePackage converts a remote package into a domain.PackageAdd. The remote package must be hosted on GitHub or GitLab.
// If the fileURL is valid and the remote package is a GitHub or GitLab repository, the remote package is fetched and
// the files are extracted. The file tree is then converted into a domain.PackageAdd. The name of the package is the name
// of the remote package. If the fileURL is not valid, an error is returned.
func (i *_importer) RemotePackage(ctx context.Context, fileURL string) (*domain.PackageAdd, error) {
	// Parse file URL
	sourceURL, err := remote.ParseRepoURL(fileURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse file URL")
	}

	// Get package information
	packageInfo, err := remote.ExtractPackageInfoFromURL(ctx, sourceURL)
	if err != nil {
		return nil, errors.Wrap(err, "extract package information")
	}

	// Determine platform
	_platform, err := platform.NewWithAuth(ctx, sourceURL.Host, nil)
	if err != nil {
		return nil, errors.Wrap(err, "identify platform")
	}

	// Get file contents
	contents, sha, err := _platform.GetFileContent(ctx, packageInfo.Repo, packageInfo.Path)
	if err != nil {
		return nil, errors.Wrap(err, "get file URL contents")
	}

	// Create package
	buf := bytes.NewBuffer(contents)
	_package, err := packageFromTOMLReader(buf)
	if err != nil {
		return nil, errors.Wrap(err, "create package")
	}

	// Set package metadata
	if _package.UpstreamURL == nil || *_package.UpstreamURL == "" {
		fileURL = sourceURL.String()
		_package.UpstreamURL = &fileURL
	}
	if _package.SHA == nil || *_package.SHA == "" {
		_package.SHA = &sha
	}

	return _package, nil
}

// RepositoryAsPackage converts a remote repository into a domain.PackageAdd. The remote repository must be a GitHub or
// GitLab repository. The remote repository is fetched and the files are extracted. The file tree is then converted
// into a domain.PackageAdd. The name of the package is the name of the remote repository.
func (i *_importer) RepositoryAsPackage(ctx context.Context, repoURL string, exclude *regexp.Regexp) (*domain.PackageAdd, error) {
	// Prepare the remote repository.
	_repoURL, err := remote.ParseRepoURL(repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse repo URL")
	}

	_platform, err := platform.NewWithAuth(ctx, _repoURL.Hostname(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "identify platform")
	}

	repoInfo, err := remote.ExtractRepoInfoFromURL(ctx, _repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "extract repo info")
	}

	// Call the remote repository.
	skipper := regexToPathSkipper(exclude)
	tree, sha, err := _platform.GetRepoTree(ctx, repoInfo, skipper)
	if err != nil {
		return nil, errors.Wrap(err, "get repo tree")
	}

	// At this point, we have the remote repository's tree. We need to convert them to templates, and then we're
	// done.
	_package := domain.NewPackageWithAutoLabel(repoInfo.Name)
	repoURL = _repoURL.String()
	_package.UpstreamURL = &repoURL
	_package.SHA = &sha
	_package.DirTree = tree

	return _package, nil
}
