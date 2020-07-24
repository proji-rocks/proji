package packageservice

import (
	"net/url"
	"path"
	"regexp"

	"github.com/nikoksr/proji/pkg/domain"

	"github.com/pkg/errors"

	"github.com/nikoksr/proji/pkg/remote"
)

// ImportFromRepoStructure imports a package from a given URL. The URL should point to a remote remote of one of the following code
// platforms: github, gitlab. Proji will imitate the structure and content of the remote and create a package
// based on it.
func (ps packageService) ImportPackageFromRepositoryStructure(url *url.URL, exclude *regexp.Regexp) (*domain.Package, error) {
	// Get code repo
	codeRepo, err := remote.NewCodeRepository(url, ps.authentication)
	if err != nil {
		return nil, errors.Wrap(err, "get code repository")
	}

	// Set package name from base name
	// E.g. https://github.com/nikoksr/proji -> proji is the base name
	name := path.Base(url.Path)
	label := pickLabel(name)
	pkg := domain.NewPackage(name, label)

	// Get templates from repo tree entries
	pkg.Templates, err = codeRepo.GetTreeEntriesAsTemplates(url, exclude)
	if err != nil {
		return nil, errors.Wrap(err, "get templates from tree entries")
	}

	// Validate package correctness
	err = isPackageValid(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "package validation")
	}
	return pkg, nil
}
