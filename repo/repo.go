package repo

import (
	"fmt"
	"net/url"
	"strings"
)

// Importer describes the behaviour of repo objects (github, gitlab).
type Importer interface {
	FilePathToRawURI(filePath string) string // Returns raw URI of a file
	LoadTreeEntries() error                  // Loads a list of tree entries of a specific repo
	Owner() string                           // Returns the name of the owner
	Repo() string                            // Returns the name of the repo
	Branch() string                          // Returns the name of the branch
	URL() *url.URL                           // Returns the url of the attached repository
}

// ParseURL parses a regular URL of a remote repository into a "cleaned up" version.
//
// Steps:
//   - Remove trailing '.git'
//   - Replace domain abbreviations with full domains
//   - Parse raw string to URL structure
//   - Make absolute if not already
func ParseURL(repoURL string) (*url.URL, error) {
	if strings.Trim(repoURL, " ") == "" {
		return nil, fmt.Errorf("can't parse empty url")
	}

	// Trim trailing '.git'
	if strings.HasSuffix(repoURL, ".git") {
		repoURL = repoURL[:len(repoURL)-len(".git")]
	}

	// domainAbbreviations defines a map of domain abbreviations like 'gh:' and their associated full domains like
	// 'https://github.com'.
	var domainAbbreviations = map[string]string{
		"gh:": "https://github.com",
		"gl:": "https://gitlab.com",
	}
	// Replace domain abbreviations like 'gh:' with the actual domain of the host
	for abbreviation, fullDomain := range domainAbbreviations {
		if strings.HasPrefix(repoURL, abbreviation) {
			repoURL = strings.Replace(repoURL, abbreviation, fullDomain, 1)
			break
		}
	}

	// Parse to URL structure
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}

	// Make absolute if not already
	if !u.IsAbs() {
		u.Scheme = "https"
		u, err = url.Parse(u.String())
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}
