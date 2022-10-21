package remote

import (
	"context"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

// Platform is an interface that defines the methods that a platform must implement. A platform is a remote repository
// that can be used to fetch and store information. For example, GitHub is a platform. We use these platforms to fetch
// and store information about packages.
type Platform interface {
	GetRepoTree(ctx context.Context, info RepoInfo, skipper PathSkipperFn) (domain.DirTree, string, error)
	GetFileContent(ctx context.Context, info RepoInfo, file string) ([]byte, string, error)
	DownloadFile(ctx context.Context, info RepoInfo, source, destination string) error
	DownloadFileRaw(ctx context.Context, source, destination string) error
	String() string
}

// PathSkipperFn is a function that returns true if the given path should be skipped.
type PathSkipperFn func(path string) (shouldSkip bool)

// DefaultPathSkipper is the default PathSkipperFn. It always returns false, meaning that no paths should be skipped.
func DefaultPathSkipper(string) bool {
	return false
}

// IsStatusCodeOK returns true if the given status code is in the 2XX range.
func IsStatusCodeOK(code int) bool {
	return code >= 200 && code < 300
}
