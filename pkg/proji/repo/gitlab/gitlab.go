package gitlab

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/tidwall/gjson"
)

// gitlab struct holds important data about a gitlab repo
type gitlab struct {
	baseURI    *url.URL
	apiBaseURL string
	ownerName  string
	repoName   string
	branchName string
}

// New creates a new gitlab repo object
func New(URL *url.URL) (repo.Importer, error) {
	if URL.Hostname() != "gitlab.com" {
		return nil, fmt.Errorf("invalid host %s", URL.Hostname())
	}

	// Parse URL
	// Examples:
	//  - /[inkscape]/[inkscape]                 	-> extracts owner and repo name; no branch name
	//  - /[inkscape]/[inkscape]/-/tree/[master]	-> extracts owner, repo and branch name
	r := regexp.MustCompile(`/([^/]+)/([^/]+)(?:/-/(?:(?:blob|tree)/([^/]+)))?`)
	specs := r.FindStringSubmatch(URL.Path)

	if specs == nil {
		return nil, fmt.Errorf("could not parse url path")
	}

	ownerName := specs[1]
	repoName := specs[2]
	branchName := specs[3]

	if ownerName == "" || repoName == "" {
		return nil, fmt.Errorf("could not extract user and/or repository name. Please check the URL")
	}

	// Default to master if no branch was defined
	if branchName == "" {
		branchName = "master"
	}

	return &gitlab{
		baseURI:    URL,
		apiBaseURL: "https://gitlab.com/api/v4/projects/",
		ownerName:  ownerName,
		repoName:   repoName,
		branchName: branchName,
	}, nil
}

// GetBaseURI returns the base URI of the repo
// You can pass the relative path to a file of that repo to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this repo.
func (g *gitlab) FilePathToRawURI(filePath string) string {
	return g.baseURI.String() +
		filepath.Join(
			"/-/raw/",
			g.branchName, "/",
			filePath,
		)
}

// GetTree gets the paths and types of the repo tree
func (g *gitlab) GetTree(filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result, error) {
	nextPage := "1"
	paths := make([]gjson.Result, 0)
	types := make([]gjson.Result, 0)
	treeReq := g.apiBaseURL + g.ownerName + "%2F" + g.repoName + "/repository/tree/?ref=" + g.branchName + "&recursive=true&per_page=100&page="

	for nextPage != "" {
		// Request repo tree
		response, err := repo.GetRequest(treeReq + nextPage)
		if err != nil {
			return nil, nil, err
		}

		// Parse the tree
		body, _ := ioutil.ReadAll(response.Body)
		treeResponse := gjson.GetMany(string(body), "#.path", "#.type")
		paths = append(paths, treeResponse[0].Array()...)
		types = append(types, treeResponse[1].Array()...)
		err = response.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		// Set next page from response header
		nextPage = response.Header.Get("X-Next-Page")
	}
	// Filter paths
	paths, types = repo.FilterPathsNTypes(paths, types, filters)
	return paths, types, nil
}

// GetOwner returns the name of the repo owner
func (g *gitlab) Owner() string { return g.ownerName }

// GetRepo returns the name of the repo
func (g *gitlab) Repo() string { return g.repoName }
