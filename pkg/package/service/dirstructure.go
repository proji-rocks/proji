package packageservice

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nikoksr/proji/internal/util"
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
)

// ImportFromFolderStructure imports a package from a given directory. Proji will imitate the
// structure and content of the directory and create a package based on it.
func (ps packageService) ImportPackageFromDirectoryStructure(path string, filters []*regexp.Regexp) (*domain.Package, error) {
	// Validate that the directory exists
	if !util.DoesPathExist(path) {
		return nil, fmt.Errorf("given directory does not exist")
	}

	// No filters given
	if filters == nil {
		filters = make([]*regexp.Regexp, 0)
	}

	// Set package name from directory base name
	name := filepath.Base(path)
	label := pickLabel(name)
	pkg := domain.NewPackage(name, label)

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
			return err
		}

		// Skip if path matches a filter
		skipPath := doesPathMatchFilter(currentPath, filters)
		if skipPath {
			return filepath.SkipDir
		}

		// Add file or folder to package
		isFile := true
		if info.IsDir() {
			isFile = false
		}
		pkg.Templates = append(pkg.Templates, &domain.Template{IsFile: isFile, Path: "", Destination: relPath})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Validate package
	err = isPackageValid(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "package validation")
	}
	return pkg, nil
}

func doesPathMatchFilter(path string, filters []*regexp.Regexp) bool {
	for _, filter := range filters {
		if filter.FindStringIndex(path) != nil {
			return true

		}
	}
	return false
}
