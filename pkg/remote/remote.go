package remote

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/remote/github"
	"github.com/nikoksr/proji/pkg/remote/gitlab"
)

type CodeRepository interface {
	GetPackageConfig(url *url.URL) (string, error)
	GetCollectionConfigs(url *url.URL, filters []*regexp.Regexp) ([]string, error)
	GetTreeEntriesAsTemplates(url *url.URL, filters []*regexp.Regexp) ([]*domain.Template, error)
}

// ParseURL parses a regular URL of a remote repository into a "cleaned up" version.
//
// Steps:
//   - Remove trailing '.git'
//   - Replace domain abbreviations with full domains
//   - Parse raw string to URL structure
//   - Make absolute if not already
func ParseURL(repoURL string) (*url.URL, error) {
	if strings.TrimSpace(repoURL) == "" {
		return nil, fmt.Errorf("can't parse empty url")
	}

	// Trim trailing '.git'
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// domainAbbreviations defines a map of domain abbreviations like 'gh:' and their associated full domains like
	// 'https://github.com'.
	var domainAbbreviations = map[string]string{
		"gh:": "https://github.com",
		"gl:": "https://gitlab.com",
	}
	// Replace domain abbreviations like 'gh:' with the actual domain of the host
	for abbreviation, fullDomain := range domainAbbreviations {
		repoURL = strings.Replace(repoURL, abbreviation, fullDomain, 1)
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

// NewCodeHostingPlatform returns the most suiting import based on the code hosting platform.
func NewCodeRepository(repoURL *url.URL, auth *config.APIAuthentication) (CodeRepository, error) {
	var platform CodeRepository
	var err error

	switch repoURL.Hostname() {
	case "github.com":
		platform, err = github.New(auth.GHToken)
	case "gitlab.com":
		platform, err = gitlab.New(auth.GLToken)
	default:
		return nil, fmt.Errorf("code hosting platform not supported")
	}
	return platform, err
}
