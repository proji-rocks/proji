package gitlab

import (
	"net/url"
	"os"
	"testing"

	"github.com/nikoksr/proji/util"
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

var apiToken = os.Getenv("PROJI_AUTH_GL_TOKEN")

// TestNew tests the creation of a new github object based on given github URLs.
func TestNew(t *testing.T) {
	util.SkipNetworkBasedTests(t)

	// These should work
	for i, repo := range goodRepos {
		g, err := New(repo.baseURI, apiToken)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		assert.Equal(t, goodRepos[i].baseURI, g.baseURI)
		assert.Equal(t, goodRepos[i].OwnerName, g.OwnerName)
		assert.Equal(t, goodRepos[i].RepoName, g.RepoName)
		assert.Equal(t, goodRepos[i].BranchName, g.BranchName)
	}

	// These should fail
	for _, URL := range badURLs {
		_, err := New(URL, apiToken)
		assert.Error(t, err)
	}
}

// TestGitLab_LoadTreeEntries tests the github method TestGetTreePathsAndTypes which tries
// to request and receive the folders paths and types of a github repo tree.
func TestGitLab_LoadTreeEntries(t *testing.T) {
	util.SkipNetworkBasedTests(t)

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
					path:      "include/util.hpp",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/util.cpp",
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
					path:      "include/util.hpp",
					entryType: "blob",
				},
				{
					id:        "77083e0ec310487f88cf875f5ea7f377ee1819ad",
					path:      "notes.txt",
					entryType: "blob",
				},
				{
					id:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
					path:      "src/util.cpp",
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
		g, err := New(test.URL, apiToken)
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

func TestGitLab_FilePathToRawURI(t *testing.T) {
	type fields struct {
		OwnerName  string
		RepoName   string
		BranchName string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test FilePathToRawURI 1",
			fields: fields{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test",
				BranchName: "master",
			},
			args: args{filePath: "/configs/test.conf"},
			want: "https://gitlab.com/nikoksr/proji-test/-/raw/master/configs/test.conf",
		},
		{
			name: "Test FilePathToRawURI 2",
			fields: fields{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test-package",
				BranchName: "develop",
			},
			args: args{filePath: "/test/some_test.go"},
			want: "https://gitlab.com/nikoksr/proji-test-package/-/raw/develop/test/some_test.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitLab{
				OwnerName:  tt.fields.OwnerName,
				RepoName:   tt.fields.RepoName,
				BranchName: tt.fields.BranchName,
			}
			if got := g.FilePathToRawURI(tt.args.filePath); got != tt.want {
				t.Errorf("FilePathToRawURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitLab_Owner(t *testing.T) {
	tests := []struct {
		name string
		got  *GitLab
		want string
	}{
		{
			name: "Test Owner 1",
			got:  goodRepos[0],
			want: "nikoksr",
		},
		{
			name: "Test Owner 2",
			got: &GitLab{
				OwnerName: "testUser247",
			},
			want: "testUser247",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.got.Owner(), "%s\n", test.name)
	}
}

func TestGitLab_Repo(t *testing.T) {
	tests := []struct {
		name string
		got  *GitLab
		want string
	}{
		{
			name: "Test Repo 1",
			got:  goodRepos[0],
			want: "proji-test",
		},
		{
			name: "Test Repo 2",
			got: &GitLab{
				RepoName: "testRepo247",
			},
			want: "testRepo247",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.got.Repo(), "%s\n", test.name)
	}
}

func TestGitLab_Branch(t *testing.T) {
	tests := []struct {
		name string
		got  *GitLab
		want string
	}{
		{
			name: "Test Branch 1",
			got:  goodRepos[0],
			want: "master",
		},
		{
			name: "Test Branch 2",
			got:  goodRepos[1],
			want: "develop",
		},
		{
			name: "Test Branch 3",
			got: &GitLab{
				BranchName: "testBranch247",
			},
			want: "testBranch247",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.got.Branch(), "%s\n", test.name)
	}
}
