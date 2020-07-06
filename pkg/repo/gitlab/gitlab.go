package gitlab

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	gl "github.com/xanzy/go-gitlab"
)

// GitLab struct holds important data about a gitlab repo
type GitLab struct {
	baseURI     *url.URL
	OwnerName   string
	RepoName    string
	BranchName  string
	TreeEntries []*gl.TreeNode
	client      *gl.Client
}

// New creates a new gitlab repo object
func New(URL *url.URL, authToken string) (*GitLab, error) {
	if URL.Hostname() != "gitlab.com" {
		return nil, fmt.Errorf("invalid host %s", URL.Hostname())
	}

	// Parse URL
	// Examples:
	//  - /[inkscape]/[inkscape]                 	-> extracts owner and repo name; no branch name
	//  - /[inkscape]/[inkscape]/-/tree/[master]	-> extracts owner, repo and branch name
	r := regexp.MustCompile(`/([^/]+)/([^/]+)(?:/(?:tree|blob)/([^/]+))?`)
	specs := r.FindStringSubmatch(URL.Path)

	if specs == nil {
		return nil, fmt.Errorf("could not parse url path")
	}

	OwnerName := specs[1]
	RepoName := specs[2]
	BranchName := specs[3]

	if OwnerName == "" || RepoName == "" {
		return nil, fmt.Errorf("could not extract user and/or repository name. Please check the URL")
	}

	// Default to master if no branch was defined
	if BranchName == "" {
		BranchName = "master"
	}

	glClient, err := gl.NewClient(authToken)
	if err != nil {
		return nil, err
	}

	return &GitLab{
		baseURI:    URL,
		OwnerName:  OwnerName,
		RepoName:   RepoName,
		BranchName: BranchName,
		client:     glClient,
	}, nil
}

// GetBaseURI returns the base URI of the repo
// You can pass the relative path to a file of that repo to receive the complete raw url for said file.
// Or you pass an empty string resulting in the base of the raw url for files of this repo.
func (g *GitLab) FilePathToRawURI(filePath string) string {
	filePath = strings.TrimPrefix(filePath, "/")
	return fmt.Sprintf("https://gitlab.com/%s/%s/-/raw/%s/%s", g.OwnerName, g.RepoName, g.BranchName, filePath)
}

// GetTreeEntries gets the paths and types of the repo tree
func (g *GitLab) LoadTreeEntries() error {
	// Reset tree entries
	g.TreeEntries = make([]*gl.TreeNode, 0)
	pid := g.OwnerName + "/" + g.RepoName
	rec := true
	nextPage := 1

	for {
		treeNodes, resp, err := g.client.Repositories.ListTree(pid, &gl.ListTreeOptions{
			Recursive: &rec,
			Ref:       &g.BranchName,
			ListOptions: gl.ListOptions{
				Page:    nextPage,
				PerPage: 100,
			},
		})
		if err != nil {
			return err
		}
		g.TreeEntries = append(g.TreeEntries, treeNodes...)

		// Break if no next page
		if resp.NextPage == 0 {
			break
		}
		nextPage = resp.NextPage
	}
	return nil
}

// Owner returns the name of the owner
func (g *GitLab) Owner() string { return g.OwnerName }

// Repo returns the name of the repo
func (g *GitLab) Repo() string { return g.RepoName }

// Repo returns the name of the branch
func (g *GitLab) Branch() string { return g.BranchName }
