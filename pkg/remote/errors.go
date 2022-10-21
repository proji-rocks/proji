package remote

import "github.com/cockroachdb/errors"

var (
	ErrPackageNotFound = errors.New("package not found")
	ErrRepoNotFound    = errors.New("repo not found")
)
