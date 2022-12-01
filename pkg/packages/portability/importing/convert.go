package importing

import (
	"encoding/json"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

// packageFromTOML unmarshals a TOML byte slice and returns a package. It returns an error if the byte slice is not a TOML
// file, or it was not able to unmarshal the TOML file.
func packageFromTOML(buf []byte) (*domain.PackageAdd, error) {
	pkg := &domain.PackageAdd{}
	if err := toml.Unmarshal(buf, pkg); err != nil {
		return nil, errors.Wrap(err, "unmarshal package")
	}

	return pkg, nil
}

// packageFromTOMLReader unmarshals a TOML reader and returns a package. It returns an error if the reader is not a TOML
// file, or it was not able to unmarshal the TOML file.
func packageFromTOMLReader(r io.Reader) (*domain.PackageAdd, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "read all")
	}

	return packageFromTOML(buf)
}

// packageFromJSON unmarshals a JSON byte slice and returns a package. It returns an error if the byte slice is not a JSON
// file, or it was not able to unmarshal the JSON file.
func packageFromJSON(buf []byte) (*domain.PackageAdd, error) {
	pkg := &domain.PackageAdd{}
	if err := json.Unmarshal(buf, pkg); err != nil {
		return nil, errors.Wrap(err, "unmarshal package")
	}

	return pkg, nil
}

// packageFromJSONReader unmarshals a JSON reader and returns a package. It returns an error if the reader is not a JSON
// file, or it was not able to unmarshal the JSON file.
func packageFromJSONReader(r io.Reader) (*domain.PackageAdd, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return &domain.PackageAdd{}, errors.Wrap(err, "read all")
	}

	return packageFromJSON(buf)
}
