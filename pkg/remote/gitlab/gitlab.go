package gitlab

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nikoksr/proji/internal/util"
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	gl "github.com/xanzy/go-gitlab"
)

const defaultTimeout = time.Second * 10

// Service struct holds important data about a gitlab remote.
type Service struct {
	isAuthenticated bool
	client          *gl.Client
	repo            *repo
}

// New creates a new gitlab remote instance.
func New(authToken string) (*Service, error) {
	s := &Service{}
	var err error
	s.client, err = gl.NewClient(authToken)
	if err != nil {
		return nil, errors.Wrap(err, "gitlab client")
	}
	return s, nil
}

func (s Service) getRepository(url *url.URL) (*repo, error) {
	if url.Hostname() != "gitlab.com" {
		return nil, fmt.Errorf("invalid host %s", url.Hostname())
	}

	// Parse URL
	// Examples:
	//  - /[inkscape]/[inkscape]                 	-> extracts owner and remote name; no branch name
	//  - /[inkscape]/[inkscape]/-/tree/[master]	-> extracts owner, remote and branch name
	regex := regexp.MustCompile(`/([^/]+)/([^/]+)(?:/(?:tree|blob)/([^/]+))?`)
	specs := regex.FindStringSubmatch(url.Path)

	if specs == nil {
		return nil, fmt.Errorf("could not parse url path")
	}

	owner := specs[1]
	repoName := specs[2]
	branch := specs[3]

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("could not extract user and/or repository name. Please check the URL")
	}

	// Default to master if no branch was defined
	if branch == "" {
		branch = "master"
	}

	return &repo{
		url:    url,
		name:   repoName,
		owner:  owner,
		branch: branch,
		client: s.client,
	}, nil
}

func (s Service) GetTreeEntriesAsTemplates(url *url.URL, filters []*regexp.Regexp) ([]*domain.Template, error) {
	repo, err := s.getRepository(url)
	if err != nil {
		return nil, errors.Wrap(err, "gitlab service")
	}

	if filters == nil {
		filters = make([]*regexp.Regexp, 0)
	}

	treeEntries, err := repo.getTreeEntries()
	if err != nil {
		return nil, errors.Wrap(err, "gitlab repository")
	}
	templates := make([]*domain.Template, 0)

	for _, entry := range treeEntries {
		skip := false
		for _, filter := range filters {
			if !skip {
				skip = true
			}
			if filter.FindStringIndex(entry.Path) != nil {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		isFile := false
		if entry.Type == "blob" {
			isFile = true
		}
		templates = append(templates, &domain.Template{
			IsFile:      isFile,
			Path:        "",
			Destination: entry.Path,
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

func (s Service) GetCollectionConfigs(url *url.URL, filters []*regexp.Regexp) ([]string, error) {
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
	return repo.downloadCollectionConfigs(treeEntries, filters)
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

func (r repo) downloadCollectionConfigs(treeEntries []*gl.TreeNode, filters []*regexp.Regexp) ([]string, error) {
	var configPaths []string
	configsBasePath := filepath.Join(os.TempDir(), "proji/configs")
	filters = append([]*regexp.Regexp{regexp.MustCompile(`configs/.*`)}, filters...)
	for _, entry := range treeEntries {
		for _, filter := range filters {
			// Skip if doesn't match filter
			if filter.FindStringIndex(entry.Path) == nil {
				continue
			}

			// Parse package url from base url
			packageURL := r.url
			packageURL.Path = filepath.Join(packageURL.Path, entry.Path)

			// Download config
			configPath, err := r.downloadPackageConfig(configsBasePath, packageURL)
			if err != nil {
				return nil, err
			}

			// Add config path to list
			configPaths = append(configPaths, configPath)
			break
		}
	}
	return configPaths, nil
}

// GetBaseURI returns the base URI of the remote
// You can pass the relative path to a file of that remote to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this remote.
func (r repo) getRawFileURL(filePath string) string {
	filePath = strings.TrimPrefix(filePath, "/")
	return fmt.Sprintf("https://gitlab.com/%s/%s/-/raw/%s/%s", r.owner, r.name, r.branch, filePath)
}
