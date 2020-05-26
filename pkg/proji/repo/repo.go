package repo

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// domainAbbreviations defines a map of domain abbreviations like 'gh:' and their associated full domains like
// 'https://github.com'.
var domainAbbreviations = map[string]string{
	"gh:": "https://github.com",
	"gl:": "https://gitlab.com",
}

// Importer describes the behaviour of repo objects (github, gitlab)
type Importer interface {
	FilePathToRawURI(filePath string) string                                  // Returns raw URI of a file
	GetTree(filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result, error) // Returns the paths and types of the repo tree
	Owner() string                                                            // Returns the name of the repo owner
	Repo() string                                                             // Returns the name of the repo
}

// GetRequest is a wrapper for the http.Get() method, handling errors and bad status codes
func GetRequest(request string) (*http.Response, error) {
	response, err := http.Get(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s returned status code %d", request, response.StatusCode)
	}
	return response, nil
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
