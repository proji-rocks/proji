package github

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/tidwall/gjson"
)

// github struct holds important data about a github repo
type github struct {
	baseURI    *url.URL
	apiBaseURL string
	ownerName  string
	repoName   string
	branchName string
	repoSHA    string
}

// setRepoSHA sets the repoSHA attribute equal to the SHA-1 of the last commit in the current branch
func (g *github) setRepoSHA() error {
	// Send request for SHA-1 of branch
	shaReq := g.apiBaseURL + g.ownerName + "/" + g.repoName + "/branches/" + g.branchName
	response, err := repo.GetRequest(shaReq)
	if err != nil {
		return err
	}

	// Parse body and try to extract SHA
	body, _ := ioutil.ReadAll(response.Body)
	repoSHA := gjson.Get(string(body), "commit.sha")
	defer response.Body.Close()
	if !repoSHA.Exists() {
		return fmt.Errorf("could not get commit sha-1 from body")
	}
	g.repoSHA = repoSHA.String()
	return nil
}

// New creates a new github repo object
func New(URL *url.URL) (repo.Importer, error) {
	if URL.Hostname() != "github.com" {
		return nil, fmt.Errorf("invalid host %s", URL.Hostname())
	}

	// Extract owner, repo and branch if given
	// Examples:
	//  - /[nikoksr]/[proji]				-> extracts owner and repo name; no branch name
	//  - /[nikoksr]/[proji]/tree/[master]	-> extracts owner, repo and branch name
	r := regexp.MustCompile(`/([^/]+)/([^/]+)(?:/tree/([^/]+))?`)
	specs := r.FindStringSubmatch(URL.Path)

	if specs == nil {
		return nil, fmt.Errorf("could not parse url")
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

	g := &github{
		baseURI:    URL,
		apiBaseURL: "https://api.github.com/repos/",
		ownerName:  ownerName,
		repoName:   repoName,
		branchName: branchName,
		repoSHA:    "",
	}
	return g, g.setRepoSHA()
}

// GetBaseURI returns the base URI of the repo
// You can pass the relative path to a file of that repo to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this repo.
func (g *github) FilePathToRawURI(filePath string) string {
	return "https://raw.githubusercontent.com/" +
		filepath.Join(
			g.ownerName,
			g.repoName,
			g.branchName,
			filePath,
		)
}

// GetTree gets the paths and types of the repo tree
func (g *github) GetTree(filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result, error) {
	// Request repo tree
	treeReq := g.apiBaseURL + g.ownerName + "/" + g.repoName + "/git/trees/" + g.repoSHA + "?recursive=1"
	response, err := repo.GetRequest(treeReq)
	if err != nil {
		return nil, nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	// Check if response was truncated
	if gjson.Get(string(body), "truncated").Bool() == true {
		return nil, nil, fmt.Errorf("the response was truncated by Github, which means that the number of items in the tree array exceeded the maximum limit.\n\nClone the repo manually with git and use 'proji class import --directory /path/to/repo' to import your local copy of that repo")
	}

	// Parse the tree
	treeResponse := gjson.GetMany(string(body), "tree.#.path", "tree.#.type")
	defer response.Body.Close()
	paths := treeResponse[0].Array()
	types := treeResponse[1].Array()

	// Filter paths
	paths, types = repo.FilterPathsNTypes(paths, types, filters)
	return paths, types, nil
}

// Owner returns the name of the repo owner
func (g *github) Owner() string { return g.ownerName }

// Repo returns the name of the repo
func (g *github) Repo() string { return g.repoName }
