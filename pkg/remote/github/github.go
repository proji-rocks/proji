package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nikoksr/proji/internal/util"

	"github.com/nikoksr/proji/pkg/domain"

	gh "github.com/google/go-github/v31/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const defaultTimeout = time.Second * 10

// Service struct holds important data about a github remote.
type Service struct {
	isAuthenticated bool
	client          *gh.Client
	repo            *repo
}

// New creates a new github remote instance.
func New(authToken string) (*Service, error) {
	s := &Service{}
	err := s.setClient(authToken)
	if err != nil {
		return nil, errors.Wrap(err, "github client")
	}
	return s, nil
}

func (s Service) getRepository(url *url.URL) (*repo, error) {
	if url.Hostname() != "github.com" {
		return nil, fmt.Errorf("invalid host %s", url.Hostname())
	}

	// Extract owner, remote and branch if given
	// Examples:
	//  - /[nikoksr]/[proji]				-> extracts owner and remote name; no branch name
	//  - /[nikoksr]/[proji]/tree/[master]	-> extracts owner, remote and branch name
	regex := regexp.MustCompile(`/([^/]+)/([^/]+)(?:/tree/([^/]+))?`)
	specs := regex.FindStringSubmatch(url.Path)

	if specs == nil {
		return nil, fmt.Errorf("could not parse url")
	}

	owner := specs[1]
	repoName := specs[2]

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("could not extract user and/or repository name")
	}

	currentRepo := &repo{
		url:    url,
		name:   repoName,
		owner:  owner,
		branch: specs[3],
		client: s.client,
	}

	err := currentRepo.setSHA(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "repository commit sha")
	}
	return currentRepo, nil
}

func (s *Service) setClient(authToken string) error {
	if len(strings.TrimSpace(authToken)) <= 0 {
		// Create an unauthenticated client
		s.client = gh.NewClient(&http.Client{Timeout: defaultTimeout})
		s.isAuthenticated = false
		return nil
	}
	// Create oauth token
	oauthToken := &oauth2.Token{AccessToken: authToken}
	if !oauthToken.Valid() {
		return fmt.Errorf("couldn't validate github access token")
	}

	// Create an authenticated client
	authenticatedClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(oauthToken))
	authenticatedClient.Timeout = defaultTimeout
	s.client = gh.NewClient(authenticatedClient)
	s.isAuthenticated = true
	return nil
}

func (s Service) GetTreeEntriesAsTemplates(url *url.URL, exclude *regexp.Regexp) ([]*domain.Template, error) {
	repo, err := s.getRepository(url)
	if err != nil {
		return nil, errors.Wrap(err, "github service")
	}

	treeEntries, err := repo.getTreeEntries()
	if err != nil {
		return nil, errors.Wrap(err, "github repository")
	}
	templates := make([]*domain.Template, 0, len(treeEntries))

	for _, entry := range treeEntries {
		path := entry.GetPath()

		// Check if exclude matches the path
		doesMatch := exclude.MatchString(path)
		if doesMatch {
			continue
		}

		// Parse to template
		isFile := false
		if entry.GetType() == "blob" {
			isFile = true
		}
		templates = append(templates, &domain.Template{
			IsFile:      isFile,
			Path:        "",
			Destination: path,
		})
	}
	return templates, nil
}

func (s Service) GetPackageConfig(url *url.URL) (string, error) {
	repo, err := s.getRepository(url)
	if err != nil {
		return "", errors.Wrap(err, "github service")
	}

	configPath := filepath.Join(os.TempDir(), "proji/configs")
	configPath, err = repo.downloadPackageConfig(configPath, url)
	if err != nil {
		return "", errors.Wrap(err, "download package config")
	}
	return configPath, nil
}

func (r repo) downloadCollectionConfigs(treeEntries []*gh.TreeEntry, exclude *regexp.Regexp) ([]string, error) {
	var configPaths []string
	configsBasePath := filepath.Join(os.TempDir(), "proji/configs")
	for _, entry := range treeEntries {
		path := entry.GetPath()
		// Check if exclude matches the path
		doesMatch := exclude.MatchString(path)
		if doesMatch {
			continue
		}

		// Parse package url from base url
		packageURL := r.url
		packageURL.Path = filepath.Join(packageURL.Path, entry.GetPath())

		// Download config
		configPath, err := r.downloadPackageConfig(configsBasePath, packageURL)
		if err != nil {
			return nil, err
		}

		// Add config path to list
		configPaths = append(configPaths, configPath)
	}
	return configPaths, nil
}

func (s Service) GetCollectionConfigs(url *url.URL, exclude *regexp.Regexp) ([]string, error) {
	// Setup repo
	repo, err := s.getRepository(url)
	if err != nil {
		return []string{}, errors.Wrap(err, "github service")
	}

	// Get tree entries from repo
	treeEntries, err := repo.getTreeEntries()
	if err != nil {
		return []string{}, errors.Wrap(err, "repo tree entries")
	}

	// Filter out only files under configs/ path
	return repo.downloadCollectionConfigs(treeEntries, exclude)
}

// GetBaseURI returns the base URI of the remote
// You can pass the relative path to a file of that remote to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this remote.
func (r repo) getRawFileURL(filePath string) string {
	filePath = strings.TrimPrefix(filePath, "/")
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", r.owner, r.name, r.branch, filePath)
}

// downloadPackageConfig downloads the config file from the given url to the given destination. Returns an error if not
// a valid download source and returns the full path of the downloaded file on success.
func (r repo) downloadPackageConfig(destination string, source *url.URL) (string, error) {
	// Extract file name
	fileName := filepath.Base(source.Path)

	// download file to temporary directory
	destination = filepath.Join(destination, fileName)

	// get raw file url ready for download
	rawSource := r.getRawFileURL(filepath.Join("configs", fileName))
	err := util.DownloadFileIfNotExists(destination, rawSource)
	if err != nil {
		return "", errors.Wrap(err, "download config file")
	}
	return destination, nil
}
