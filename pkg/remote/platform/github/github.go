package github

import (
	"context"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	gh "github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/httputil"
	"github.com/nikoksr/proji/pkg/logging"
	"github.com/nikoksr/proji/pkg/remote"
)

// Compile-time check to ensure that GitHub implements the remote.Platform interface.
var _ remote.Platform = &GitHub{}

// GitHub is an implementation of the platform.Platform interface.
type GitHub struct {
	client *gh.Client
}

const defaultTimeout = 10 * time.Second

// httpClient is the default HTTP client used by the GitHub client.
var httpClient = &http.Client{
	Timeout: defaultTimeout,
}

func newWithAuth(ctx context.Context, token string) *GitHub {
	logger := logging.FromContext(ctx)

	if token == "" {
		logger.Debugf("no token provided, using anonymous GitHub client")

		return &GitHub{client: gh.NewClient(httpClient)}
	}

	logger.Debugf("token provided, using authenticated GitHub client")

	// Token given, create a new client with the token.
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// Create http and GitHub client.
	_httpClient := oauth2.NewClient(ctx, tokenSource)
	_httpClient.Timeout = defaultTimeout

	client := gh.NewClient(_httpClient)

	return &GitHub{client: client}
}

// newWithContext creates a new GitHub remote instance.
func newWithContext(ctx context.Context, token string) *GitHub {
	return newWithAuth(ctx, token)
}

// New creates a new GitHub remote instance.
func New(ctx context.Context, token string) *GitHub {
	return newWithContext(ctx, token)
}

// GetRepoTree returns the list of entries for the given repository as a domain.DirTree. It also returns the SHA of the
// repo tree. This is useful for versioning packages.
func (g *GitHub) GetRepoTree(ctx context.Context, info remote.RepoInfo, skip remote.PathSkipperFn) (domain.DirTree, string, error) {
	logger := logging.FromContext(ctx)

	if skip == nil {
		logger.Debugf("using default path skipper")
		skip = remote.DefaultPathSkipper
	}

	if info.Owner == "" || info.Name == "" || info.Ref == "" {
		return nil, "", errors.New("missing owner, name or ref")
	}

	logger.Debugf("getting repo tree for %s/%s@%s", info.Owner, info.Name, info.Ref)
	repoTree, resp, err := g.client.Git.GetTree(ctx, info.Owner, info.Name, info.Ref, true)

	// Special case for when the repo is not found. Instead of the messy api error message, we return a more friendly
	// error.
	if resp.StatusCode == http.StatusNotFound {
		return nil, "", remote.ErrRepoNotFound
	}

	// Now check if there were any other errors.
	if err != nil {
		return nil, "", errors.Wrap(err, "get tree")
	}

	// And finally make a sanity check on the response code.
	if !remote.IsStatusCodeOK(resp.StatusCode) {
		return nil, "", errors.New("get tree responded with code " + resp.Status)
	}

	logger.Debugf("processing repo tree for %s/%s@%s", info.Owner, info.Name, info.Ref)
	tree := make(domain.DirTree, 0, len(repoTree.Entries))

	for _, entry := range repoTree.Entries {
		if entry == nil || entry.Path == nil || entry.Type == nil {
			continue
		}

		path := entry.GetPath()
		if skip(path) {
			continue
		}

		tree = append(tree, &domain.DirEntry{
			Path:  path,
			IsDir: entry.GetType() == "tree",
		})
	}

	return tree, repoTree.GetSHA(), nil
}

func (g *GitHub) getContent(ctx context.Context, info remote.RepoInfo, path string) (*gh.RepositoryContent, string, error) {
	if info.Owner == "" || info.Name == "" || info.Ref == "" {
		return nil, "", errors.New("missing owner, name or ref")
	}

	opts := &gh.RepositoryContentGetOptions{
		Ref: info.Ref,
	}

	contents, _, resp, err := g.client.Repositories.GetContents(ctx, info.Owner, info.Name, path, opts)

	// Special case for when the package config is not found. Instead of the messy api error message, we return a more
	// friendly error.
	if resp.StatusCode == http.StatusNotFound {
		return nil, "", remote.ErrPackageNotFound
	}

	// Now check if there were any other errors.
	if err != nil {
		return nil, "", errors.Wrap(err, "get file contents")
	}

	// And finally make a sanity check on the response code.
	if !remote.IsStatusCodeOK(resp.StatusCode) {
		return nil, "", errors.New("get file contents responded with code " + resp.Status)
	}

	if contents == nil {
		return nil, "", errors.New("file is empty")
	}

	return contents, contents.GetSHA(), nil
}

func (g *GitHub) getFileContent(ctx context.Context, info remote.RepoInfo, file string) (*gh.RepositoryContent, string, error) {
	content, sha, err := g.getContent(ctx, info, file)
	if err != nil {
		return nil, "", err
	}

	if content.GetType() != "file" {
		return nil, "", errors.Newf("content is of type %q, not a file", content.GetType())
	}

	return content, sha, nil
}

// GetFileContent returns the content of the given file. If the file is a directory, an error will be returned. If the
// file does not exist, an error will be returned.
func (g *GitHub) GetFileContent(ctx context.Context, info remote.RepoInfo, file string) ([]byte, string, error) {
	logger := logging.FromContext(ctx)

	logger.Debugf("getting content for file %s/%s@%s:%s", info.Owner, info.Name, info.Ref, file)
	content, sha, err := g.getFileContent(ctx, info, file)
	if err != nil {
		return nil, "", err
	}

	logger.Debugf("decoding content for file %s/%s@%s:%s", info.Owner, info.Name, info.Ref, file)
	decodedContent, err := content.GetContent()
	if err != nil {
		return nil, "", errors.Wrap(err, "decode file content")
	}

	return []byte(decodedContent), sha, nil
}

// toRawContentURL returns the URL to the raw content of the given file. Info as well as the file path are required.
func toRawContentURL(info remote.RepoInfo, file string) (string, error) {
	if info.Owner == "" || info.Name == "" || info.Ref == "" || file == "" {
		return "", errors.New("missing owner, name, ref or file")
	}

	return "https://raw.githubusercontent.com/" + info.Owner + "/" + info.Name + "/" + info.Ref + "/" + file, nil
}

// DownloadFileRaw downloads the given file from the remote repository. Internally, it uses the raw content URL to fetch
// the file contents. This is done to avoid using up the GitHub API rate limit. The given URL is required and must be
// a valid raw content URL pointing to the file.
func (g *GitHub) DownloadFileRaw(ctx context.Context, fileURL, dest string) error {
	return httputil.DownloadFile(ctx, fileURL, dest)
}

// DownloadFile downloads the given file from the remote repository. Internally, it uses the raw content URL to fetch
// the file contents. This is done to avoid using up the GitHub API rate limit.
// The given file path is required and relative to the root of the repository. The full download path will be created
// automatically.
func (g *GitHub) DownloadFile(ctx context.Context, info remote.RepoInfo, file, dest string) error {
	url, err := toRawContentURL(info, file)
	if err != nil {
		return errors.Wrap(err, "to raw content url")
	}

	return g.DownloadFileRaw(ctx, url, dest)
}

// String returns the string representation of the platform. It is used to identify the platform and satisfy the
// fmt.Stringer interface.
func (g *GitHub) String() string {
	return "github"
}
