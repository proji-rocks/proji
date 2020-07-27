package packageservice

import (
	"net/url"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/pkg/errors"
)

// ImportFromRepoStructure imports a package from a given URL. The URL should point to a remote remote of one of the following code
// platforms: github, gitlab. Proji will imitate the structure and content of the remote and create a package
// based on it.
func (ps packageService) ImportPackageFromRemote(url *url.URL) (*domain.Package, error) {
	// Get code repo
	codeRepo, err := remote.NewCodeRepository(url, ps.authentication)
	if err != nil {
		return nil, errors.Wrap(err, "get code repository")
	}

	// Download package config and get path of file
	configFile, err := codeRepo.GetPackageConfig(url)
	if err != nil {
		return nil, errors.Wrap(err, "get package config")
	}

	// Import package from that config file
	pkg, err := ps.ImportPackageFromConfig(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "import from config")
	}

	// Validate package correctness
	err = isPackageValid(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "package validation")
	}
	return pkg, nil
}
