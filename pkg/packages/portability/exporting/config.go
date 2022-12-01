package exporting

import (
	"context"
	"encoding/json"
	"os"

	"github.com/nikoksr/proji/pkg/packages/portability"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

func write(_ context.Context, file *os.File, data []byte) error {
	// Write data to file.
	written, err := file.Write(data)
	if err != nil {
		return err
	}

	// Check if all data was written.
	if written != len(data) {
		return errors.New("incomplete write")
	}
	if written == 0 {
		return errors.New("no data written")
	}

	return nil
}

func toConfig(ctx context.Context, pkg *domain.PackageExport, dir, fileType string) (string, error) {
	if pkg == nil {
		return "", errors.New("package is nil")
	}

	var err error
	var data []byte

	fileName := "proji-" + pkg.Name + ".*." + fileType

	switch fileType {
	case portability.FileTypeTOML:
		data, err = toml.Marshal(pkg)
	case portability.FileTypeJSON:
		data, err = json.MarshalIndent(pkg, "", "  ")
	default:
		err = portability.ErrUnsupportedConfigFileType
	}
	if err != nil {
		return "", err
	}

	// Open file; if dir is empty, a temporary file will be created.
	file, err := os.CreateTemp(dir, fileName)
	if err != nil {
		return "", errors.Wrap(err, "create temporary file")
	}
	defer func() { _ = file.Close() }()

	err = write(ctx, file, data)
	if err != nil {
		return "", errors.Wrap(err, "write package config")
	}

	// On success, return the file path.
	return file.Name(), nil
}

// ToConfig writes the given package to the given destination directory. If the destination is empty, a temporary file
// will be created. The caller is responsible for deleting the file. If the destination is not empty, the file will be
// overwritten.
func (e *_exporter) ToConfig(ctx context.Context, pkg *domain.PackageExport, dir, fileType string) (string, error) {
	return toConfig(ctx, pkg, dir, fileType)
}
