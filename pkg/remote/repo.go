package remote

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/nikoksr/simplog"
)

const (
	HostGitHub = "github.com"
	HostGitLab = "gitlab.com"
)

type (
	// RepoInfo is a struct that contains the most basic information about a repository. It is used to transfer info
	// between functions.
	RepoInfo struct {
		Owner string
		Name  string
		Ref   string
	}

	// PackageInfo is a struct that contains the most basic information about a package. It is used to transfer info
	// between functions. It holds a RepoInfo and a path to the actual package config file.
	PackageInfo struct {
		Repo RepoInfo
		Path string
	}
)

// domainAbbreviations defines a map of domain abbreviations like 'gh:' and their associated full domains like
// 'https://github.com'.
var domainAbbreviations = map[string]string{
	"gh:": "https://github.com/",
	"gl:": "https://gitlab.com/",
}

// ParseRepoURL parses a regular string URL of a remote repository into a "cleaned up" version. This function on purpose
// does not check if a platform is supported. This should really not be part of its logic. This function should turn
// only turn a string URL into a valid repo URL no matter what the platform is. The business logic about what platforms
// are supported should be handled by the caller or through functions like platform.New, which takes a hostname and
// returns a platform or an error.
//
// Steps:
//   - Remove trailing '.git'
//   - Replace domain abbreviations with full domains
//   - Parse raw string to URL structure and thus validate it
//   - Turn trailing @branch into branch name
func ParseRepoURL(repoURL string) (*url.URL, error) {
	repoURL = strings.TrimSpace(repoURL)

	if repoURL == "" {
		return nil, errors.New("empty URL")
	}

	// Trim trailing '.git'
	repoURL = strings.TrimSuffix(repoURL, "/")
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Replace domain abbreviations like 'gh:' with the actual domain of the host
	for abbreviation, fullDomain := range domainAbbreviations {
		if !strings.HasPrefix(repoURL, abbreviation) {
			continue
		}

		repoURL = strings.TrimPrefix(repoURL, abbreviation)
		repoURL = strings.TrimPrefix(repoURL, "/")
		repoURL = fullDomain + repoURL
	}

	// Split path by trailing @<ref>
	urlParts := strings.Split(repoURL, "@")
	if len(urlParts) > 2 {
		return nil, errors.New("invalid repository URL")
	}

	repoURL = urlParts[0]

	var branch string
	if len(urlParts) == 2 {
		branch = urlParts[1]
	}

	// General URL validation
	_repoURL, err := url.Parse(repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse repo URL")
	}

	// Check if platform is supported; if supported, check if a branch was specified and if so, turn it into a platform
	// specific branch name.
	switch _repoURL.Host {
	case HostGitHub:
		if branch != "" {
			_repoURL.Path += "/tree/" + branch
		}
	case HostGitLab:
		if branch != "" {
			_repoURL.Path += "/-/tree/" + branch
		}
	}

	return _repoURL, nil
}

const defaultBranch = "main"

var (
	regexGitHubRepo = regexp.MustCompile(`^/?([^/]+)/([^/]+)(?:/(?:tree|blob)/([^/]+)(?:/(.+))?)?`)
	regexGitLabRepo = regexp.MustCompile(`^/?([^/]+)/([^/]+)(?:(?:/-)?/(?:tree|blob)/([^/]+)(?:/(.+))?)?`)
)

// extractRepoInfoFromURL extracts the owner, name and branch from the given URL. If the URL is invalid, an error will
// be returned instead. The only supported platforms are GitHub and GitLab.
//
// GitHub examples:
//   - [nikoksr]/[proji]				                                    -> extracts owner and remote name
//   - [nikoksr]/[proji]/tree/[main]	                                    -> extracts owner, remote and branch name
//   - [nikoksr]/[proji]/tree/[fd24446df4766b987c5be0a79dd11c7bebd5dbd5]    -> extracts owner, remote and commit sha
//
// GitLab examples:
//   - [nikoksr]/[proji]				                                      -> extracts owner and remote name
//   - [nikoksr]/[proji]/-/tree/[main]	                                      -> extracts owner, remote and branch name
//   - [nikoksr]/[proji]/-/tree/[fd24446df4766b987c5be0a79dd11c7bebd5dbd5]    -> extracts owner, remote and commit sha
func extractInfoFromURL(ctx context.Context, sourceURL *url.URL) (info RepoInfo, packageConf string, err error) {
	logger := simplog.FromContext(ctx)

	if sourceURL == nil {
		return RepoInfo{}, "", errors.New("empty URL")
	}

	// Determine correct regex to use based on the hostname.
	var re *regexp.Regexp
	switch sourceURL.Host {
	case HostGitHub:
		re = regexGitHubRepo
	case HostGitLab:
		re = regexGitLabRepo
	default:
		return RepoInfo{}, "", errors.Errorf("unsupported platform: %s", sourceURL.Host)
	}

	// Try to extract meta information from the URL.
	logger.Debugf("extracting meta information from URL: %s", sourceURL.String())
	matches := re.FindStringSubmatch(sourceURL.Path)
	if matches == nil || len(matches) < 4 {
		return RepoInfo{}, "", errors.New("invalid repository URL")
	}

	// Validate extracted repo information.
	owner := matches[1]
	name := matches[2]
	if owner == "" || name == "" {
		return RepoInfo{}, "", errors.New("missing owner or repo name")
	}

	// Extract branch information and set sane default if not found.
	branch := matches[3]
	if branch == "" {
		logger.Debugf("no ref found in URL, using default: %s", defaultBranch)
		branch = defaultBranch
	}

	// Extract potential file path.
	path := ""
	if len(matches) > 4 {
		path = matches[4]
	}

	logger.Debugf("extracted meta information: owner=%s, name=%s, branch=%s, path=%s", owner, name, branch, path)

	return RepoInfo{
		Owner: owner,
		Name:  name,
		Ref:   branch,
	}, path, nil
}

// ExtractRepoInfoFromURL extracts the owner, name and branch from the given URL. If the URL is invalid, an error will
// be returned instead. The only supported platforms are GitHub and GitLab.
//
// GitHub examples:
//   - [nikoksr]/[proji]/tree/[main]/[configs/test.toml]    -> extracts owner, remote. branch name and config file path
//
// GitLab examples:
//   - [nikoksr]/[proji]/-/tree/[main]/[configs/test.toml]	-> extracts owner, remote, branch name and config file path
func ExtractRepoInfoFromURL(ctx context.Context, repoURL *url.URL) (RepoInfo, error) {
	info, _, err := extractInfoFromURL(ctx, repoURL)

	return info, err
}

// ExtractPackageInfoFromURL extracts the owner, name, branch and package config path from the given URL. If the URL is
// invalid, an error will be returned instead. The only supported platforms are GitHub and GitLab.
//
// GitHub examples:
//   - [nikoksr]/[proji]				                                    -> extracts owner and remote name
//   - [nikoksr]/[proji]/tree/[main]	                                    -> extracts owner, remote and branch name
//   - [nikoksr]/[proji]/tree/[fd24446df4766b987c5be0a79dd11c7bebd5dbd5]    -> extracts owner, remote and commit sha
//
// GitLab examples:
//   - [nikoksr]/[proji]				                                      -> extracts owner and remote name
//   - [nikoksr]/[proji]/-/tree/[main]	                                      -> extracts owner, remote and branch name
//   - [nikoksr]/[proji]/-/tree/[fd24446df4766b987c5be0a79dd11c7bebd5dbd5]    -> extracts owner, remote and commit sha
func ExtractPackageInfoFromURL(ctx context.Context, repoURL *url.URL) (PackageInfo, error) {
	info, path, err := extractInfoFromURL(ctx, repoURL)

	return PackageInfo{
		Repo: info,
		Path: path,
	}, err
}
