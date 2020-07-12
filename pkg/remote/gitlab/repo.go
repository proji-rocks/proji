package gitlab

import (
	"net/url"

	gl "github.com/xanzy/go-gitlab"
)

type repo struct {
	url    *url.URL
	name   string
	owner  string
	branch string
	client *gl.Client
}

// getRepoTreeEntries gets the paths and types of the remote tree.
func (r repo) getTreeEntries() ([]*gl.TreeNode, error) {
	// Reset tree entries
	treeEntries := make([]*gl.TreeNode, 0)
	pid := r.owner + "/" + r.name
	recursive := true
	nextPage := 1

	listTreeOptions := &gl.ListTreeOptions{
		Recursive: &recursive,
		Ref:       &r.branch,
		ListOptions: gl.ListOptions{
			Page:    nextPage,
			PerPage: 100,
		},
	}

	for {
		treeNodes, resp, err := r.client.Repositories.ListTree(pid, listTreeOptions)
		if err != nil {
			return nil, err
		}
		treeEntries = append(treeEntries, treeNodes...)

		// Break if no next page
		if resp.NextPage == 0 {
			break
		}
		nextPage = resp.NextPage
	}
	return treeEntries, nil
}
