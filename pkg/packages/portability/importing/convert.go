package importing

import (
	"encoding/json"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

//
// JSON
//

func packageFromJSON(buf []byte) (*domain.PackageAdd, error) {
	pkg := &domain.PackageAdd{}
	err := json.Unmarshal(buf, &pkg)

	return pkg, errors.Wrap(err, "unmarshal json")
}

func packageFromJSONReader(r io.Reader) (*domain.PackageAdd, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "read all")
	}

	return packageFromJSON(buf)
}

func packagesFromJSON(buf []byte) ([]domain.Package, error) {
	var packages []domain.Package
	err := json.Unmarshal(buf, &packages)

	return packages, errors.Wrap(err, "unmarshal json")
}

func PackagesFromJSONReader(r io.Reader) ([]domain.Package, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}

	return packagesFromJSON(buf)
}

//
// TOML
//

func packageFromTOML(buf []byte) (*domain.PackageAdd, error) {
	data, err := toml.LoadBytes(buf)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "load toml")
	}

	pkg := &domain.PackageAdd{}
	err = data.Unmarshal(pkg)

	return pkg, errors.Wrap(err, "unmarshal toml")
}

func packageFromTOMLReader(r io.Reader) (*domain.PackageAdd, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "read all")
	}

	return packageFromTOML(buf)
}
