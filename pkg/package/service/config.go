package packageservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func (ps packageService) ImportPackageFromConfig(path string) (*domain.Package, error) {
	// Validate file path
	isJson, err := checkConfigPath(path)
	if err != nil {
		return nil, err
	}

	pkg := domain.NewPackage("", "")
	if !isJson {
		err = ps.unmarshalToml(path, pkg)
	} else {
		err = ps.unmarshalJson(path, pkg)
	}
	if err != nil {
		return pkg, err
	}

	// Validate package
	err = isPackageValid(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "package validation")
	}
	return pkg, nil
}

func (ps packageService) ExportPackageToConfig(pkg domain.Package, destination string, json bool) (string, error) {
	fileExtension := ".toml"
	if json {
		fileExtension = ".json"
	}
	confName := filepath.Join(destination, "proji-"+pkg.Name+fileExtension)
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()

	if !json {
		return confName, ps.getTomlEncoderFunction(conf)(pkg)
	}
	return confName, ps.getJsonEncoderFunction(conf)(pkg)
}

func (ps packageService) getTomlEncoderFunction(conf *os.File) func(v interface{}) error {
	return toml.NewEncoder(conf).Order(toml.OrderPreserve).Encode
}

func (ps packageService) getJsonEncoderFunction(conf *os.File) func(v interface{}) error {
	return json.NewEncoder(conf).Encode
}

func (ps packageService) unmarshalToml(path string, pkg *domain.Package) error {
	// Load file
	file, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal config into package
	err = file.Unmarshal(pkg)
	return err
}

func (ps packageService) unmarshalJson(path string, pkg *domain.Package) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, pkg)
	return err
}

func checkConfigPath(path string) (json bool, err error) {
	isJson := strings.HasSuffix(path, ".json")
	// Check if it is a toml file
	if !strings.HasSuffix(path, ".toml") && !isJson {
		return false, fmt.Errorf("config file has to be of type 'toml' or 'json'")
	}

	// Check if file is empty
	conf, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "config file info")
	}
	if conf.Size() == 0 {
		return false, fmt.Errorf("config file is empty")
	}
	return isJson, nil
}
