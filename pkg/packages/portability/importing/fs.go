package importing

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nikoksr/proji/pkg/packages/portability"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

// LocalPackage reads a TOML file and returns a package. It returns an error if the file is not a TOML file, or it was not
// able to unmarshal the TOML file.
func (i *_importer) LocalPackage(_ context.Context, path string) (_package *domain.PackageAdd, err error) {
	// Load config
	file, err := os.Open(path)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "open config file")
	}
	defer func() { _ = file.Close() }()

	// Read config; detect file extension and parse accordingly. If the file is not a TOML file, return an error.
	ext := filepath.Ext(path)
	switch ext {
	case ".toml":
		_package, err = packageFromTOMLReader(file)
	case ".json":
		_package, err = packageFromJSONReader(file)
	default:
		return &domain.PackageAdd{}, portability.ErrUnsupportedConfigFileType
	}

	return _package, errors.Wrap(err, "unmarshal config file")
}

// LocalFolderAsPackage reads a directory and returns a package. It returns an error if the directory does not exist,
// or it was not able to packageFromTOMLReader the directory.
func (i *_importer) LocalFolderAsPackage(_ context.Context, path string, exclude *regexp.Regexp) (*domain.PackageAdd, error) {
	// Pick package name and label
	name := filepath.Base(path)
	_package := domain.NewPackageWithAutoLabel(name)

	hasExcludePattern := exclude != nil && exclude.String() != ""

	// Scan directory
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip base directory
		if path == currentPath {
			return nil
		}

		// Extract relative path
		relPath, err := filepath.Rel(path, currentPath)
		if err != nil {
			return errors.Wrap(err, "extract relative path")
		}

		// Skip paths that match exclude pattern. Just continue if there is no exclude pattern in the first place.
		if hasExcludePattern && exclude.MatchString(relPath) {
			return filepath.SkipDir
		}

		// Append file or folder as template to package
		_package.DirTree.Entries = append(_package.DirTree.Entries, &domain.DirEntry{
			Path:  relPath,
			IsDir: info.IsDir(),
		})

		return nil
	})

	return _package, errors.Wrap(err, "walk directory")
}
