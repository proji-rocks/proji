package packageservice

import (
	"fmt"
	"net/url"
	"regexp"
	"sync"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/pkg/errors"
)

func (ps packageService) ImportPackagesFromCollection(url *url.URL, filters []*regexp.Regexp) ([]*domain.Package, error) {
	// Get code repo
	codeRepo, err := remote.NewCodeRepository(url, ps.authentication)
	if err != nil {
		return nil, errors.Wrap(err, "get code repository")
	}

	// Download package configs collection and get their path
	configFiles, err := codeRepo.GetCollectionConfigs(url, filters)
	if err != nil {
		return nil, errors.Wrap(err, "get package configs")
	}

	// Import all packages
	numConfigFiles := len(configFiles)
	packages := make([]*domain.Package, 0, numConfigFiles)
	var wg sync.WaitGroup
	wg.Add(numConfigFiles)
	errs := make(chan error, numConfigFiles)

	for _, configFile := range configFiles {
		go func(configFile string) {
			defer wg.Done()

			// Import package from config file
			pkg, err := ps.ImportPackageFromConfig(configFile)
			if err != nil {
				errs <- errors.Wrap(err, "import from config")
				return
			}

			// Validate package correctness
			err = isPackageValid(pkg)
			if err != nil {
				errs <- errors.Wrap(err, "package validation")
				return
			}
			packages = append(packages, pkg)
		}(configFile)
	}

	wg.Wait()
	close(errs)

	var errMsg string
	err = nil
	for e := range errs {
		if e != nil {
			errMsg += fmt.Sprintf("%v\n", e)
		}
	}

	if len(errMsg) > 0 {
		err = errors.New(errMsg)
	}

	return packages, err
}
