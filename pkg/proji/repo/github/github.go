package github

import (
	"fmt"
	"io/ioutil"

	"github.com/nikoksr/proji/pkg/proji/repo"

	"github.com/tidwall/gjson"
)

// github struct holds important data about a github repo
type github struct {
	apiBaseURI string
	userName   string
	repoName   string
	branchName string
	repoSHA    string
}

// New creates a new github repo object
func New(userName, repoName string) (repo.Importer, error) {
	g := &github{apiBaseURI: "https://api.github.com/repos/", userName: userName, repoName: repoName, repoSHA: ""}
	return g, g.setRepoSHA()
}

// GetUserName returns the name of the repo owner
func (g *github) GetUserName() string { return g.userName }

// GetRepoName returns the name of the repo
func (g *github) GetRepoName() string { return g.repoName }

// GetBranchName returns the branch name
func (g *github) GetBranchName() string { return g.branchName }

// setRepoSHA sets the repoSHA attribute equal to the SHA-1 of the last commit in the current branch
func (g *github) setRepoSHA() error {
	// Send request for SHA-1 of branch
	shaReq := g.apiBaseURI + g.userName + "/" + g.repoName + "/branches/master"
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

// GetTreePathsAndTypes gets the paths and types of the repo tree
func (g *github) GetTreePathsAndTypes() ([]gjson.Result, []gjson.Result, error) {
	// Request repo tree
	treeReq := g.apiBaseURI + g.userName + "/" + g.repoName + "/git/trees/" + g.repoSHA + "?recursive=1"
	response, err := repo.GetRequest(treeReq)
	if err != nil {
		return nil, nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	// Check if response was truncated
	if gjson.Get(string(body), "truncated").Bool() == true {
		return nil, nil, fmt.Errorf("the response was truncated by Github, which means that the number of items in the tree array exceeded the maximum limit.\n\nClone the repo manually with git and use 'proji class import --directory /path/to/repo' to import the local instance of that repo")
	}

	// Parse the tree
	treeResponse := gjson.GetMany(string(body), "tree.#.path", "tree.#.type")
	defer response.Body.Close()
	paths := treeResponse[0].Array()
	types := treeResponse[1].Array()
	return paths, types, nil
}
