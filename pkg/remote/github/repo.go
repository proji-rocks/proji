package github

import (
	"context"
	"fmt"
	"net/url"

	gh "github.com/google/go-github/v31/github"
)

type repo struct {
	url    *url.URL
	name   string
	owner  string
	branch string
	sha    string
	client *gh.Client
}

// setRepoSHA sets the repoSHA attribute equal to the SHA-1 of the last commit in the current branch.
func (r *repo) setSHA(ctx context.Context) error {
	if r.branch == "" {
		/*
			r := &gh.Repository{}
			r, _, err := g.client.Repositories.Get(ctx, g.OwnerName, g.RepoName)
			if err != nil {
				return err
			}
			g.BranchName = r.GetDefaultBranch()
		*/
		// Default to master branch for now. The above uses too many API calls and Github's API limit gets exceeded
		// too quickly.
		r.branch = "master"
	}

	branch, _, err := r.client.Repositories.GetBranch(ctx, r.owner, r.name, r.branch)
	if err != nil {
		return err
	}
	r.sha = branch.GetCommit().GetSHA()
	return nil
}

// getRepoTreeEntries gets the paths and types of the remote tree.
func (r repo) getTreeEntries() ([]*gh.TreeEntry, error) {
	tree, _, err := r.client.Git.GetTree(context.Background(), r.owner, r.name, r.sha, true)
	if err != nil {
		return nil, err
	}
	if tree.GetTruncated() {
		return nil, fmt.Errorf("the response was truncated by Github, which means that the number of items in the tree array exceeded the maximum limit.\n\nClone the remote manually with git and use 'proji package import --directory /path/to/remote' to import the local instance of that remote")
	}
	return tree.Entries, nil
}
