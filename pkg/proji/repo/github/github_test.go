package github

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	gh "github.com/google/go-github/v31/github"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/stretchr/testify/assert"
)

var goodRepos = []*GitHub{
	{
		baseURI:    &url.URL{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test"},
		OwnerName:  "nikoksr",
		RepoName:   "proji-test",
		BranchName: "master",
		repoSHA:    "b4fc28f09ac57e314d27e9b9133b1ebc03bec2f1",
	},
	{
		baseURI:    &url.URL{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/tree/develop"},
		OwnerName:  "nikoksr",
		RepoName:   "proji-test",
		BranchName: "develop",
		repoSHA:    "f07d0b57cd6b468b331be03699f15faf4f9dd910",
	},
}

var badURLs = []*url.URL{
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/does-not-exist"},
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/tree/dead-branch"},
	{Scheme: "https", Host: "github.com", Path: ""},
	{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/blaa"},
	{Scheme: "https", Host: "google.com", Path: ""},
}

// TestNew tests the creation of a new github object based on given github URLs.
func TestNew(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	// These should work
	for i, repo := range goodRepos {
		g, err := New(repo.baseURI)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		assert.Equal(t, goodRepos[i].baseURI, g.baseURI)
		assert.Equal(t, goodRepos[i].OwnerName, g.OwnerName)
		assert.Equal(t, goodRepos[i].RepoName, g.RepoName)
		assert.Equal(t, goodRepos[i].BranchName, g.BranchName)
		assert.Equal(t, goodRepos[i].repoSHA, g.repoSHA)
	}

	// These should fail
	for _, URL := range badURLs {
		_, err := New(URL)
		assert.Error(t, err)
	}
}

// TestLoadTreeEntries tests the github method TestGetTreePathsAndTypes which tries
// to request and receive the folders paths and types of a github repo tree.
func TestGitHub_LoadTreeEntries(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	type testEntry struct {
		sha       string
		path      string
		entryType string
		URL       string
	}

	tests := []struct {
		URL         *url.URL
		wantErr     bool
		treeEntries []*testEntry
	}{
		{
			URL:     goodRepos[0].baseURI,
			wantErr: false,
			treeEntries: []*testEntry{
				{
					sha:       "ce67480c3cd24e7dd675a7486233231c050f2c2e",
					path:      ".vscode",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/ce67480c3cd24e7dd675a7486233231c050f2c2e",
				},
				{
					sha:       "5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
					path:      ".vscode/c_cpp_properties.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
				},
				{
					sha:       "cf646956cf7745868f005a3b0fc622fa0390b3d7",
					path:      ".vscode/launch.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/cf646956cf7745868f005a3b0fc622fa0390b3d7",
				},
				{
					sha:       "ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
					path:      ".vscode/tasks.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
				},
				{
					sha:       "a16196bf1875a1054b731e47c528bdfc828c0649",
					path:      "CMakeLists.txt",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/a16196bf1875a1054b731e47c528bdfc828c0649",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "README.md",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
					path:      "include",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "include/helper.hpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "f22f80dfb366d311404859100709fcc348668aff",
					path:      "src",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/f22f80dfb366d311404859100709fcc348668aff",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/helper.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "b3cf51681c44016f9234f67dbd00ee49704b0021",
					path:      "src/main.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/b3cf51681c44016f9234f67dbd00ee49704b0021",
				},
				{
					sha:       "c7ae9824dc95ed736692e3a5a55bbe15ddfe250f",
					path:      "test",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/c7ae9824dc95ed736692e3a5a55bbe15ddfe250f",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "test/testHelper.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
			},
		},
		{
			URL:     goodRepos[1].baseURI,
			wantErr: false,
			treeEntries: []*testEntry{
				{
					sha:       "ce67480c3cd24e7dd675a7486233231c050f2c2e",
					path:      ".vscode",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/ce67480c3cd24e7dd675a7486233231c050f2c2e",
				},
				{
					sha:       "5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
					path:      ".vscode/c_cpp_properties.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
				},
				{
					sha:       "cf646956cf7745868f005a3b0fc622fa0390b3d7",
					path:      ".vscode/launch.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/cf646956cf7745868f005a3b0fc622fa0390b3d7",
				},
				{
					sha:       "ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
					path:      ".vscode/tasks.json",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
				},
				{
					sha:       "a16196bf1875a1054b731e47c528bdfc828c0649",
					path:      "CMakeLists.txt",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/a16196bf1875a1054b731e47c528bdfc828c0649",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "README.md",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
					path:      "include",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "include/helper.hpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "77083e0ec310487f88cf875f5ea7f377ee1819ad",
					path:      "notes.txt",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/77083e0ec310487f88cf875f5ea7f377ee1819ad",
				},
				{
					sha:       "f22f80dfb366d311404859100709fcc348668aff",
					path:      "src",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/f22f80dfb366d311404859100709fcc348668aff",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/helper.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
				{
					sha:       "b3cf51681c44016f9234f67dbd00ee49704b0021",
					path:      "src/main.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/b3cf51681c44016f9234f67dbd00ee49704b0021",
				},
				{
					sha:       "c7ae9824dc95ed736692e3a5a55bbe15ddfe250f",
					path:      "test",
					entryType: "tree",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/trees/c7ae9824dc95ed736692e3a5a55bbe15ddfe250f",
				},
				{
					sha:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "test/testHelper.cpp",
					entryType: "blob",
					URL:       "https://api.github.com/repos/nikoksr/proji-test/git/blobs/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
				},
			},
		},
		{
			URL:         badURLs[0],
			wantErr:     true,
			treeEntries: nil,
		},
	}

	for _, test := range tests {
		g, err := New(test.URL)
		assert.Equal(t, err != nil, test.wantErr, "LoadTreeEntries() error = %v, wantErr %v", err, test.wantErr)

		if g == nil {
			continue
		}

		err = g.LoadTreeEntries()
		assert.Equal(t, err != nil, test.wantErr, "LoadTreeEntries() error = %v, wantErr %v", err, test.wantErr)

		for i, entry := range g.TreeEntries {
			assert.Equal(t, entry.GetSHA(), test.treeEntries[i].sha)
			assert.Equal(t, entry.GetPath(), test.treeEntries[i].path)
			assert.Equal(t, entry.GetType(), test.treeEntries[i].entryType)
		}
	}
}

func TestGitHub_setRepoSHA(t *testing.T) {
	httpClient := http.Client{Timeout: 10 * time.Second}
	tests := []struct {
		name    string
		g       *GitHub
		wantErr bool
		want    string
	}{
		{
			name: "Valid repo SHA 1",
			g: &GitHub{
				OwnerName:   goodRepos[0].OwnerName,
				RepoName:    goodRepos[0].RepoName,
				BranchName:  goodRepos[0].BranchName,
				repoSHA:     "",
				baseURI:     goodRepos[0].baseURI,
				TreeEntries: make([]*gh.TreeEntry, 0),
				client:      gh.NewClient(&httpClient),
			},
			wantErr: false,
			want:    goodRepos[0].repoSHA,
		},
		{
			name: "Valid repo SHA 2",
			g: &GitHub{
				OwnerName:   goodRepos[1].OwnerName,
				RepoName:    goodRepos[1].RepoName,
				BranchName:  goodRepos[1].BranchName,
				repoSHA:     "",
				baseURI:     goodRepos[0].baseURI,
				TreeEntries: make([]*gh.TreeEntry, 0),
				client:      gh.NewClient(&httpClient),
			},
			wantErr: false,
			want:    goodRepos[1].repoSHA,
		},
		{
			name: "Invalid repo",
			g: &GitHub{
				OwnerName:   "nikoksr",
				RepoName:    "does-not-exist",
				BranchName:  "master",
				repoSHA:     "",
				baseURI:     badURLs[0],
				TreeEntries: make([]*gh.TreeEntry, 0),
				client:      gh.NewClient(&httpClient),
			},
			wantErr: true,
			want:    "",
		},
	}

	for _, test := range tests {
		err := test.g.setRepoSHA(context.Background())
		assert.Equal(t, test.wantErr, err != nil)
		if err == nil {
			assert.Equal(t, test.want, test.g.repoSHA)
		}
	}
}

func TestGetBranchName(t *testing.T) {
	tests := []struct {
		name string
		got  *GitHub
		want string
	}{
		{
			name: "Test Owner 1",
			got:  goodRepos[0],
			want: "nikoksr",
		},
		{
			name: "Test Owner 2",
			got: &GitHub{
				OwnerName: "testUser247",
			},
			want: "testUser247",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.got.Owner(), "%s\n", test.name)
	}
}

func TestGitHub_Repo(t *testing.T) {
	tests := []struct {
		name string
		got  *GitHub
		want string
	}{
		{
			name: "Test Repo 1",
			got:  goodRepos[0],
			want: "proji-test",
		},
		{
			name: "Test Repo 2",
			got: &GitHub{
				RepoName: "testRepo247",
			},
			want: "testRepo247",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.got.Repo(), "%s\n", test.name)
	}
}

func TestGetownerName(t *testing.T) {
	tests := []struct {
		name string
		g    *github
		want string
	}{
		{
			name: "",
			g: &github{
				apiBaseURL: ghAPIBase,
				ownerName:  "nikoksr",
				repoName:   "proji",
				branchName: "master",
				repoSHA:    "",
			},
			want: "nikoksr",
		},
		{
			name: "",
			g: &github{
				apiBaseURL: ghAPIBase,
				ownerName:  "golang",
				repoName:   "go",
				branchName: "master",
				repoSHA:    "",
			},
			want: "golang",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.g.ownerName, "%s\n", test.name)
	}
}
