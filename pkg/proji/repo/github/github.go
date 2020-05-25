package github

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/tidwall/gjson"
)

// github struct holds important data about a github repo
type github struct {
	baseURI    string
	apiBaseURI string
	userName   string
	repoName   string
	branchName string
	repoSHA    string
}

// setRepoSHA sets the repoSHA attribute equal to the SHA-1 of the last commit in the current branch
func (g *github) setRepoSHA() error {
	// Send request for SHA-1 of branch
	shaReq := g.apiBaseURI + g.userName + "/" + g.repoName + "/branches/" + g.branchName
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
func New(URL string) (repo.Importer, error) {
	// Parse URL
	// Examples:
	//  - [https://github.com/[nikoksr]/[proji]]                -> extracts base uri, user and repo name; no branch name
	//  - [https://github.com/[nikoksr]/[proji]]/tree/[master]  -> extracts base uri, user, repo and branch name
	r := regexp.MustCompile(`((?:(?:http|https)://)?github.com/([^/]+)/([^/]+))(?:/tree/([^/]+))?`)
	specs := r.FindStringSubmatch(URL)

	if specs == nil || len(specs) < 5 {
		return nil, fmt.Errorf("could not parse url")
	}

	baseURI := specs[1]
	userName := specs[2]
	repoName := specs[3]
	branchName := specs[4]

	if userName == "" || repoName == "" {
		return nil, fmt.Errorf("could not extract user and/or repository name. Please check the URL")
	}

	// Default to master if no branch was defined
	if branchName == "" {
		branchName = "master"
	}

	g := &github{
		baseURI:    baseURI,
		apiBaseURI: "https://api.github.com/repos/",
		userName:   userName,
		repoName:   repoName,
		branchName: branchName,
		repoSHA:    "",
	}
	return g, g.setRepoSHA()
}

// GetBaseURI returns the base URI of the repo
func (g *github) GetBaseURI() string { return g.baseURI }

// GetBaseURI returns the base URI of the repo
// You can pass the relative path to a file of that repo to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this repo.
func (g *github) GetRawURI(filePath string) string {
	return "https://raw.githubusercontent.com/" +
		filepath.Join(
			g.userName,
			g.repoName,
			g.branchName,
			filePath,
		)
}

// GetUserName returns the name of the repo owner
func (g *github) GetUserName() string { return g.userName }

// GetRepoName returns the name of the repo
func (g *github) GetRepoName() string { return g.repoName }

// GetBranchName returns the branch name
func (g *github) GetBranchName() string { return g.branchName }

// GetTree gets the paths and types of the repo tree
func (g *github) GetTree(filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result, error) {
	// Request repo tree
	treeReq := g.apiBaseURI + g.userName + "/" + g.repoName + "/git/trees/" + g.repoSHA + "?recursive=1"
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
