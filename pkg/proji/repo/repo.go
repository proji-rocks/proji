package repo

import (
	"fmt"
	"net/url"
	"strings"
)

// domainAbbreviations defines a map of domain abbreviations like 'gh:' and their associated full domains like
// 'https://github.com'.
var domainAbbreviations = map[string]string{
	"gh:": "https://github.com",
	"gl:": "https://gitlab.com",
}

// Importer describes the behaviour of repo objects (github, gitlab)
type Importer interface {
	FilePathToRawURI(filePath string) string // Returns raw URI of a file
	LoadTreeEntries() error                  // Loads a list of tree entries of a specific repo
	Owner() string                           // Returns the name of the owner
	Repo() string                            // Returns the name of the repo
	Branch() string                          // Returns the name of the branch
}

// ParseURL parses a regular URL of a remote repository into a "cleaned up" version.
//
// Steps:
//   - Remove trailing '.git'
//   - Replace domain abbreviations with full domains
//   - Parse raw string to URL structure
//   - Make absolute if not already
func ParseURL(URL string) (*url.URL, error) {
	// Trim trailing '.git'
	if strings.HasSuffix(URL, ".git") {
		URL = URL[:len(URL)-len(".git")]
	}

	// Replace domain abbreviations like 'gh:' with the actual domain of the host
	for abbreviation, fullDomain := range domainAbbreviations {
		if strings.HasPrefix(URL, abbreviation) {
			URL = strings.Replace(URL, abbreviation, fullDomain, 1)
			break
		}
	}

	// Parse to URL structure
	u, err := url.Parse(URL)
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

// FilterPathsNTypes converts the git types blob and tree to proji types folder and file
func FilterPathsNTypes(paths, types []gjson.Result, filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result) {
	if filters == nil {
		return paths, types
	}

	filteredPaths := make([]gjson.Result, 0)
	filteredTypes := make([]gjson.Result, 0)

	for idx, path := range paths {
		for _, filter := range filters {
			if filter.FindStringIndex(path.String()) != nil {
				filteredPaths = append(filteredPaths, path)
				filteredTypes = append(filteredTypes, types[idx])
				break
			}
		}
	}
	return filteredPaths, filteredTypes
}
