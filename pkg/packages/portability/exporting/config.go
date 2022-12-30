package exporting

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/packages/portability"
)

func encodeJSON(data *bytes.Buffer, pkg *domain.PackageConfig) error {
	enc := json.NewEncoder(data)
	enc.SetIndent("", "  ")

	return enc.Encode(pkg)
}

func encodeTOML(data *bytes.Buffer, pkg *domain.PackageConfig) error {
	enc := toml.NewEncoder(data)
	enc.SetIndentTables(true)
	enc.SetIndentSymbol("  ")

	return enc.Encode(pkg)
}

func write(_ context.Context, file *os.File, data *bytes.Buffer) error {
	// If no data is given, we return nil. This is an okay behavior, because the caller might want to write an empty
	// file.
	if data == nil {
		return nil
	}

	// If no file but data is given, we return an error. This is not okay, because data that's intended to be written
	// was given.
	if file == nil {
		return errors.New("file is nil")
	}

	// Write data to file.
	bufferSize := data.Len() // Before writing, get the buffer size.
	written, err := data.WriteTo(file)
	if err != nil {
		return err
	}

	// Check if all data was written.
	if written != int64(bufferSize) {
		return errors.Newf("incomplete write; written %d bytes, expected %d bytes", written, bufferSize)
	}

	return nil
}

func toConfig(ctx context.Context, pkg *domain.PackageConfig, dir, fileType string) (string, error) {
	if pkg == nil {
		return "", errors.New("package is nil")
	}

	var err error
	data := new(bytes.Buffer)
	fileName := "proji-" + pkg.Name + ".*." + fileType

	switch fileType {
	case portability.FileTypeTOML:
		err = encodeTOML(data, pkg)
	case portability.FileTypeJSON:
		err = encodeJSON(data, pkg)
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
func (e *_exporter) ToConfig(ctx context.Context, pkg *domain.PackageConfig, dir, fileType string) (string, error) {
	return toConfig(ctx, pkg, dir, fileType)
}
