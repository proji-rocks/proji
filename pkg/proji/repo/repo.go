package repo

import (
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"
)

// Importer describes the behaviour of repo objects (github, gitlab)
type Importer interface {
	FilePathToRawURI(filePath string) string                                  // Returns raw URI of a file
	GetTree(filters []*regexp.Regexp) ([]gjson.Result, []gjson.Result, error) // Returns the paths and types of the repo tree
	Owner() string                                                            // Returns the name of the repo owner
	Repo() string                                                             // Returns the name of the repo
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
