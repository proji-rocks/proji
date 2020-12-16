package github

import (
	"context"
	"fmt"
	"net/http"
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

func (r *repo) setBranch(ctx context.Context, isAuthenticated bool) error {
	// Skip if a branch was already set
	if r.branch != "" {
		return nil
	}

	// Only do this if the user is authenticated to github otherwise this would fill up your rate limit to quickly.
	if isAuthenticated {
		// Get an instance of our github repo directly from the github api library
		ghRepo, response, err := r.client.Repositories.Get(ctx, r.owner, r.name)
		if err != nil {
			return err
		}

		// Validate the response
		err = gh.CheckResponse(response.Response)
		if err != nil {
			return err
		}

		// Try to get the repo's default branch
		r.branch = ghRepo.GetDefaultBranch()
		if r.branch == "" {
			return fmt.Errorf("failed to load the default branch for the repository")
		}
	}

	// Just try the two default branches - master being the old and main the new default branch of new repos on gh.
	// Yes, I know it would probably be more efficient at this point in time to test the master branch first since
	// many older repos use that as their default, but since I support github's change of the default branch name and
	// the distribution will shift more towards the main branch in the future anyway, we'll test the main branch first
	// here. We're not wasting any api calls here.
	branches := []string{"main", "master"}
	for _, branch := range branches {
		doesExist, err := doesBranchExist(r.owner, r.name, branch)
		if err != nil {
			return err
		}
		if doesExist {
			r.branch = branch
			break
		}
	}

	return nil
}

// doesBranchExist checks if the branch of a repository exists by checking the http get return code of the repo url
// including the branch name.
func doesBranchExist(owner, repo, branch string) (bool, error) {
	repoURL := fmt.Sprintf("https://github.com/%s/%s/tree/%s", owner, repo, branch)
	resp, err := http.Get(repoURL)
	if err != nil {
		return false, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return true, nil
	}

	return false, nil
}

// setRepoSHA sets the repoSHA attribute equal to the SHA-1 of the last commit in the current branch.
func (r *repo) setSHA(ctx context.Context) error {
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
