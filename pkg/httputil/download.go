package httputil

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/logging"
	"github.com/nikoksr/proji/pkg/remote"
)

// DownloadFile downloads the given source and writes it to the given destination.
func DownloadFile(ctx context.Context, source, destination string) error {
	logger := logging.FromContext(ctx)

	// Prepare the destination directory.
	logger.Debugf("preparing destination directory: %s", path.Dir(destination))
	err := os.MkdirAll(path.Dir(destination), 0o755)
	if err != nil {
		return errors.Wrap(err, "prepare destination directory")
	}

	// Create the file; use .tmp extension, so we don't overwrite a valid file if there is an error.
	logger.Debugf("creating temporary destination file: %s.tmp", destination)
	file, err := os.Create(destination + ".tmp")
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer func() { _ = file.Close() }()

	// Download the file's contents.
	logger.Debugf("downloading file: %s", source)
	resp, err := Get(ctx, source)
	if err != nil {
		return errors.Wrap(err, "get file")
	}
	defer func() { _ = resp.Body.Close() }()

	if !remote.IsStatusCodeOK(resp.StatusCode) {
		return errors.New("download responded with code " + resp.Status)
	}
	if resp.ContentLength == 0 {
		return errors.New("download responded with empty body")
	}

	// Before writing contents to the file, make sure the context was not canceled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Write the body to file.
	logger.Debugf("saving downloaded contents of %s to %s", source, destination)
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "write file")
	}
	if written != resp.ContentLength {
		return errors.New("content length mismatch")
	}
	if written == 0 {
		return errors.New("no content written")
	}

	// Download was successful. Remove the .tmp extension.
	logger.Debugf("download complete; renaming temporary file to %s", destination)
	err = os.Rename(destination+".tmp", destination)

	return errors.Wrap(err, "rename file")
}
