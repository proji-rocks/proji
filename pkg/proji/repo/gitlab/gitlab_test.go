package gitlab

import (
	"net/url"
	"testing"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/stretchr/testify/assert"
)

var goodRepos = []*GitLab{
	{
		baseURI:    &url.URL{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test"},
		OwnerName:  "nikoksr",
		RepoName:   "proji-test",
		BranchName: "master",
	},
	{
		baseURI:    &url.URL{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test/tree/develop"},
		OwnerName:  "nikoksr",
		RepoName:   "proji-test",
		BranchName: "develop",
	},
}

var badURLs = []*url.URL{
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/does-not-exist"},
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/-/tree/dead-branch"},
	{Scheme: "https", Host: "github.com", Path: ""},
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
	}

	// These should fail
	for _, URL := range badURLs {
		_, err := New(URL)
		assert.Error(t, err)
	}
}

// TestGitLab_LoadTreeEntries tests the github method TestGetTreePathsAndTypes which tries
// to request and receive the folders paths and types of a github repo tree.
func TestGitLab_LoadTreeEntries(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	type testEntry struct {
		id        string
		path      string
		entryType string
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
					id:        "ce67480c3cd24e7dd675a7486233231c050f2c2e",
					path:      ".vscode",
					entryType: "tree",
				},
				{
					id:        "7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
					path:      "include",
					entryType: "tree",
				},
				{
					id:        "f22f80dfb366d311404859100709fcc348668aff",
					path:      "src",
					entryType: "tree",
				},
				{
					id:        "1c0a24c11ed67d83dec1cc26d252ab3d52da9f3f",
					path:      "test",
					entryType: "tree",
				},
				{
					id:        "5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
					path:      ".vscode/c_cpp_properties.json",
					entryType: "blob",
				},
				{
					id:        "cf646956cf7745868f005a3b0fc622fa0390b3d7",
					path:      ".vscode/launch.json",
					entryType: "blob",
				},
				{
					id:        "ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
					path:      ".vscode/tasks.json",
					entryType: "blob",
				},
				{
					id:        "a16196bf1875a1054b731e47c528bdfc828c0649",
					path:      "CMakeLists.txt",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "README.md",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "include/helper.hpp",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/helper.cpp",
					entryType: "blob",
				},
				{
					id:        "b3cf51681c44016f9234f67dbd00ee49704b0021",
					path:      "src/main.cpp",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "test/TestHelper.cpp",
					entryType: "blob",
				},
			},
		},
		{
			URL:     goodRepos[1].baseURI,
			wantErr: false,
			treeEntries: []*testEntry{
				{
					id:        "ce67480c3cd24e7dd675a7486233231c050f2c2e",
					path:      ".vscode",
					entryType: "tree",
				},
				{
					id:        "7213500b0fd381eb9c8e57cfdfd9b0387bcabce0",
					path:      "include",
					entryType: "tree",
				},
				{
					id:        "f22f80dfb366d311404859100709fcc348668aff",
					path:      "src",
					entryType: "tree",
				},
				{
					id:        "1c0a24c11ed67d83dec1cc26d252ab3d52da9f3f",
					path:      "test",
					entryType: "tree",
				},
				{
					id:        "5de84ef9d7019f8b47493e5d111dc1d60cf7a452",
					path:      ".vscode/c_cpp_properties.json",
					entryType: "blob",
				},
				{
					id:        "cf646956cf7745868f005a3b0fc622fa0390b3d7",
					path:      ".vscode/launch.json",
					entryType: "blob",
				},
				{
					id:        "ecbd3b5084f7657eea227f09e8fe5c0972d98d0b",
					path:      ".vscode/tasks.json",
					entryType: "blob",
				},
				{
					id:        "a16196bf1875a1054b731e47c528bdfc828c0649",
					path:      "CMakeLists.txt",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "README.md",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "include/helper.hpp",
					entryType: "blob",
				},
				{
					id:        "77083e0ec310487f88cf875f5ea7f377ee1819ad",
					path:      "notes.txt",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/helper.cpp",
					entryType: "blob",
				},
				{
					id:        "b3cf51681c44016f9234f67dbd00ee49704b0021",
					path:      "src/main.cpp",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "test/TestHelper.cpp",
					entryType: "blob",
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
			assert.Equal(t, entry.ID, test.treeEntries[i].id)
			assert.Equal(t, entry.Path, test.treeEntries[i].path)
			assert.Equal(t, entry.Type, test.treeEntries[i].entryType)
		}
	}
}

	// These should fail
	var badRepoObjects = []repo.Importer{
		&gitlab{
			apiBaseURL: "",
			ownerName:  "",
			repoName:   "",
			branchName: "",
		},
		&gitlab{
			apiBaseURL: glAPIBase,
			ownerName:  "nikoksr",
			repoName:   "proji-test",
			branchName: "does_not_exist",
		},
	}

	for _, ghRepo := range badRepoObjects {
		paths, types, err := ghRepo.GetTree(nil)
		assert.Error(t, err)
		assert.Nil(t, paths)
		assert.Nil(t, types)
	}
}

func TestGetBranchName(t *testing.T) {
	tests := []struct {
		name string
		g    *gitlab
		want string
	}{
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "nikoksr",
				repoName:   "proji-test",
				branchName: "master",
			},
			want: "master",
		},
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "nikoksr",
				repoName:   "proji-test",
				branchName: "develop",
			},
			want: "develop",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.g.branchName, "%s\n", test.name)
	}
}

func TestGetRepoName(t *testing.T) {
	tests := []struct {
		name string
		g    *gitlab
		want string
	}{
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "nikoksr",
				repoName:   "proji-test",
				branchName: "master",
			},
			want: "proji-test",
		},
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "inkscape",
				repoName:   "inkscape",
				branchName: "develop",
			},
			want: "inkscape",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.g.repoName, "%s\n", test.name)
	}
}

func TestGetownerName(t *testing.T) {
	tests := []struct {
		name string
		g    *gitlab
		want string
	}{
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "nikoksr",
				repoName:   "proji-test",
				branchName: "master",
			},
			want: "nikoksr",
		},
		{
			name: "",
			g: &gitlab{
				apiBaseURL: glAPIBase,
				ownerName:  "inkscape",
				repoName:   "inkscape",
				branchName: "master",
			},
			want: "inkscape",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.g.ownerName, "%s\n", test.name)
	}
}
