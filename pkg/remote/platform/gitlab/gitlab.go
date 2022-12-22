package gitlab

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	gl "github.com/xanzy/go-gitlab"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/httputil"
	"github.com/nikoksr/proji/pkg/remote"
)

// Compile-time check to ensure GitLab implements the remote.Platform interface.
var _ remote.Platform = &GitLab{}

// GitLab represents a GitLab remote instance. It's used to interact with the GitLab API.
type GitLab struct {
	client *gl.Client
}

// httpClient is the default HTTP client used by the GitLab client.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// NewWithOAuth creates a new GitLab remote instance. It uses the given OAuth2 token to authenticate with the GitLab
// API.
func NewWithOAuth(ctx context.Context, token string) (*GitLab, error) {
	logger := simplog.FromContext(ctx)

	logger.Debugf("creating new GitLab client with OAuth")
	client, err := gl.NewOAuthClient(token, gl.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &GitLab{client: client}, nil
}

// NewWithPAT creates a new GitLab remote instance. It uses the given PAT to authenticate with the GitLab API.
func NewWithPAT(ctx context.Context, token string) (*GitLab, error) {
	logger := simplog.FromContext(ctx)

	logger.Debugf("creating new GitLab client with PAT")
	client, err := gl.NewClient(token, gl.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &GitLab{client: client}, nil
}

// New creates a new GitLab remote instance. It uses the given token to authenticate with the GitLab API. It tries to
// guess the authentication method based on the token. If the token begins with "oauth2:" or "Bearer ", it will use the
// OAuth2 authentication method. Otherwise, it will use the Personal Access Token authentication method.
//
// If you want to avoid this guessing, use NewWithOAuth or NewWithPAT, respectively.
func New(ctx context.Context, token string) (*GitLab, error) {
	if strings.HasPrefix(token, "oauth2:") || strings.HasPrefix(token, "Bearer:") {
		return NewWithOAuth(ctx, token)
	}

	return NewWithPAT(ctx, token)
}

func (g *GitLab) getRepoSHA(info remote.RepoInfo) (string, error) {
	commit, _, err := g.client.Commits.GetCommit(info.Owner+"/"+info.Name, info.Ref, nil)

	sha := ""
	if commit != nil {
		sha = commit.ID
	}

	return sha, err
}

// GetRepoTree returns the list of entries for the given repository as a domain.DirTree.
func (g *GitLab) GetRepoTree(ctx context.Context, info remote.RepoInfo, skip remote.PathSkipperFn) (*domain.DirTree, string, error) {
	logger := simplog.FromContext(ctx)

	pid := info.Owner + "/" + info.Name
	recursive := true
	listOptions := &gl.ListTreeOptions{
		Recursive: &recursive,
		Ref:       &info.Ref,
		ListOptions: gl.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	// Get latest commit SHA
	logger.Debugf("getting latest commit SHA for %s@%s", pid, info.Ref)
	sha, err := g.getRepoSHA(info)
	if err != nil {
		logger.Debugf("failed to get repo SHA for %s: %s", pid, err)
		sha = ""
	}
	logger.Debugf("setting sha of %s to %s", pid, sha)

	var tree domain.DirTree

	for {
		// Get the next page of tree entries.
		logger.Debugf("getting page %d with size %d for %s@%s", listOptions.Page, listOptions.PerPage, pid, info.Ref)
		nodes, resp, err := g.client.Repositories.ListTree(pid, listOptions)
		if err != nil {
			return nil, "", errors.Wrap(err, "list tree")
		}
		if !remote.IsStatusCodeOK(resp.StatusCode) {
			return nil, "", errors.New("list tree responded with code " + resp.Status)
		}

		// Store the tree entries to the list.
		if tree.Entries == nil {
			tree.Entries = make([]*domain.DirEntry, 0, resp.TotalItems)
		}
		for _, node := range nodes {
			if node == nil || skip(node.Path) {
				continue
			}

			tree.Entries = append(tree.Entries, &domain.DirEntry{
				Path:  node.Path,
				IsDir: node.Type == "tree",
			})
		}

		// If there are no more pages, we're done.
		if resp.NextPage == 0 {
			logger.Debugf("no more pages for %s@%s", pid, info.Ref)
			break
		}

		// Otherwise, fetch the next page.
		listOptions.ListOptions.Page = resp.NextPage
	}

	return &tree, sha, nil
}

// GetFileContent returns the content of the given file. If the file is a directory, an error will be returned. If the
// file does not exist, an error will be returned.
func (g *GitLab) GetFileContent(ctx context.Context, info remote.RepoInfo, file string) ([]byte, string, error) {
	logger := simplog.FromContext(ctx)

	if info.Owner == "" || info.Name == "" || info.Ref == "" {
		return nil, "", errors.New("missing owner, name or ref")
	}

	pid := info.Owner + "/" + info.Name
	opts := &gl.GetFileOptions{
		Ref: &info.Ref,
	}

	logger.Debugf("getting file %s@%s:%s", pid, info.Ref, file)
	fileContents, resp, err := g.client.RepositoryFiles.GetFile(pid, file, opts, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "get file")
	}

	if !remote.IsStatusCodeOK(resp.StatusCode) {
		return nil, "", errors.New("get file responded with code " + resp.Status)
	}

	if fileContents == nil {
		return nil, "", errors.New("get file returned nil")
	}

	var content []byte
	switch fileContents.Encoding {
	case "base64":
		logger.Debugf("decoding base64 file content %s@%s:%s", pid, info.Ref, file)
		content, err = base64.StdEncoding.DecodeString(fileContents.Content)
	case "":
		logger.Debugf("content of %s@%s:%s is not encoded", pid, info.Ref, file)
		content = []byte(fileContents.Content)
	default:
		err = errors.Newf("unsupported encoding: %s", fileContents.Encoding)
	}
	if err != nil {
		return nil, "", errors.Wrap(err, "decode file contents")
	}

	return content, fileContents.SHA256, nil
}

// toRawContentURL returns the URL to the raw content of the given file. Info as well as the file path are required.
func toRawContentURL(info remote.RepoInfo, file string) (string, error) {
	if info.Owner == "" || info.Name == "" || info.Ref == "" || file == "" {
		return "", errors.New("missing owner, name, ref or file")
	}

	return "https://gitlab.com/" + info.Owner + "/" + info.Name + "/raw/" + info.Ref + "/" + file, nil
}

// DownloadFileRaw downloads the given file from the remote repository. Internally, it uses the raw content URL to fetch
// the file contents. This is done to avoid using up the GitLab API rate limit. The given URL is required and must be
// a valid raw content URL pointing to the file.
func (g *GitLab) DownloadFileRaw(ctx context.Context, fileURL, dest string) error {
	return httputil.DownloadFile(ctx, fileURL, dest)
}

// DownloadFile downloads the given file from the remote repository. Internally, it uses the raw content URL to fetch
// the file contents. This is done to avoid using up the GitLab API rate limit.
// The given file path is required and relative to the root of the repository. The full download path will be created
// automatically.
func (g *GitLab) DownloadFile(ctx context.Context, info remote.RepoInfo, file, dest string) error {
	url, err := toRawContentURL(info, file)
	if err != nil {
		return errors.Wrap(err, "to raw content url")
	}

	return g.DownloadFileRaw(ctx, url, dest)
}

// String returns the string representation of the platform. It is used to identify the platform and satisfy the
// fmt.Stringer interface.
func (g *GitLab) String() string {
	return "gitlab"
}
