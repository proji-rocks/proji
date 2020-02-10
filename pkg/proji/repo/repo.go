package repo

import (
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"
)

// Importer describes the behaviour of repo objects (github, gitlab)
type Importer interface {
	GetUserName() string                                           // Returns the name of the repo owner
	GetRepoName() string                                           // Returns the name of the repo
	GetBranchName() string                                         // Returns the branch name
	GetTreePathsAndTypes() ([]gjson.Result, []gjson.Result, error) // Get the paths and types of the repo tree
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
