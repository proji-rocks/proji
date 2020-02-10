package gitlab

import (
	"io/ioutil"

	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/tidwall/gjson"
)

// gitlab struct holds important data about a gitlab repo
type gitlab struct {
	apiBaseURI string
	userName   string
	repoName   string
}

// New creates a new gitlab repo object
func New(userName, repoName string) repo.Importer {
	return &gitlab{apiBaseURI: "https://gitlab.com/api/v4/projects/", userName: userName, repoName: repoName}
}

// GetTreePathsAndTypes gets the paths and types of the repo tree
func (g *gitlab) GetTreePathsAndTypes() ([]gjson.Result, []gjson.Result, error) {
	nextPage := "1"
	paths := make([]gjson.Result, 0)
	types := make([]gjson.Result, 0)
	treeReq := g.apiBaseURI + g.userName + "%2F" + g.repoName + "/repository/tree/?recursive=true&per_page=100&page="

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
