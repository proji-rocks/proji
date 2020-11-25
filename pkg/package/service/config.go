package packageservice

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func (ps packageService) ImportPackageFromConfig(path string) (*domain.Package, error) {
	// Validate file path
	err := isConfigPathValid(path)
	if err != nil {
		return nil, err
	}

	// Load file
	file, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal config into package
	pkg := domain.NewPackage("", "")
	err = file.Unmarshal(pkg)
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

func (ps packageService) ExportPackageToConfig(pkg domain.Package, destination string) (string, error) {
	confName := filepath.Join(destination, "proji-"+pkg.Name+".toml")
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()

	// return confName, toml.NewEncoder(conf).Order(toml.OrderPreserve).Encode(pkg)
	return confName, toml.NewEncoder(conf).Encode(pkg)
}
}

func isConfigPathValid(path string) error {
	// Check if it is a toml file
	if !strings.HasSuffix(path, ".toml") {
		return fmt.Errorf("config file has to be of type 'toml'")
	}

	// Check if file is empty
	conf, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "config file info")
	}
	if conf.Size() == 0 {
		return fmt.Errorf("config file is empty")
	}
	return nil
}
