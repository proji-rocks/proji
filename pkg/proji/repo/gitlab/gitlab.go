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
	//  - https://gitlab.com/[inkscape]/[inkscape]                  -> extracts user and repo name; no branch name
	//  - https://gitlab.com/[inkscape]/[inkscape]/-/tree/[master]  -> extracts user, repo and branch name
	r := regexp.MustCompile(`gitlab.com/(?P<User>[^/]+)/(?P<Repo>[^/]+)(/-/tree/(?P<Branch>[^/]+))?`)
	specs := r.FindStringSubmatch(repoURLPath)

	if specs == nil || len(specs) < 5 {
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

// GetUserName returns the name of the repo owner
func (g *gitlab) GetUserName() string { return g.userName }

// GetRepoName returns the name of the repo
func (g *gitlab) GetRepoName() string { return g.repoName }

// GetBranchName returns the branch name
func (g *gitlab) GetBranchName() string { return g.branchName }

// GetTreePathsAndTypes gets the paths and types of the repo tree
func (g *gitlab) GetTreePathsAndTypes() ([]gjson.Result, []gjson.Result, error) {
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
	return paths, types, nil
}
