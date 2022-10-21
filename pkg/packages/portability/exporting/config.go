package exporting

import (
	"context"
	"encoding/json"
	"os"

	"github.com/cockroachdb/errors"

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

func toConfig(ctx context.Context, pkg *domain.Package, dir string) (string, error) {
	if pkg == nil {
		return "", errors.New("package is nil")
	}

	// Open file; if dir is empty, a temporary file will be created.
	fileName := "proji-" + pkg.Name + ".*.json"
	file, err := os.CreateTemp(dir, fileName)
	if err != nil {
		return "", errors.Wrap(err, "create temporary file")
	}
	defer func() { _ = file.Close() }()

	// Write package to file.
	pkgJSON, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "marshal package")
	}

	err = write(ctx, file, pkgJSON)
	if err != nil {
		return "", errors.Wrap(err, "write package config")
	}

	// On success, return the file path.
	return file.Name(), nil
}

// ToConfig writes the given package to the given destination directory. If the destination is empty, a temporary file
// will be created. The caller is responsible for deleting the file. If the destination is not empty, the file will be
// overwritten.
func (e *_exporter) ToConfig(ctx context.Context, pkg *domain.Package, dir string) (string, error) {
	return toConfig(ctx, pkg, dir)
}
