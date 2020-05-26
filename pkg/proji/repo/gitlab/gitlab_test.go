package gitlab

import (
	"net/url"
	"testing"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const glAPIBase = "https://gitlab.com/api/v4/projects/"

var goodURLs = []*url.URL{
	{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test"},
	{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test/-/tree/develop"},
}

var goodRepoObjects = []repo.Importer{
	&gitlab{
		baseURI:    goodURLs[0],
		apiBaseURL: glAPIBase,
		ownerName:  "nikoksr",
		repoName:   "proji-test",
		branchName: "master",
	},
	&gitlab{
		baseURI:    goodURLs[1],
		apiBaseURL: glAPIBase,
		ownerName:  "nikoksr",
		repoName:   "proji-test",
		branchName: "develop",
	},
}

// TestNew tests the creation of a new github object based on given github URLs.
func TestNew(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	// These should work
	for i, URL := range goodURLs {
		g, err := New(URL)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		assert.Equal(t, goodRepoObjects[i], g)
	}

	// These should fail
	var badURLs = []*url.URL{
		{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/does-not-exist"},
		{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test/-/tree/dead-branch"},
		{Scheme: "https", Host: "gitlab.com", Path: ""},
		{Scheme: "https", Host: "github.com", Path: "/nikoksr/blaa"},
		{Scheme: "https", Host: "google.com", Path: ""},
	}

	for _, URL := range badURLs {
		glRepo, err := New(URL)
		if err == nil {
			_, _, err = glRepo.GetTree(nil)
		}
		assert.Error(t, err)
	}
}

// TestGetTreePathsAndTypes tests the github method TestGetTreePathsAndTypes which tries
// to request and receive the folders paths and types of a github repo tree.
func TestGetTreePathsAndTypes(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	var goodRepoPathsResults = [][]gjson.Result{
		{
			{Type: gjson.Type(3), Raw: "\".vscode\"", Str: ".vscode", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include\"", Str: "include", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src\"", Str: "src", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test\"", Str: "test", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/c_cpp_properties.json\"", Str: ".vscode/c_cpp_properties.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/launch.json\"", Str: ".vscode/launch.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/tasks.json\"", Str: ".vscode/tasks.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"CMakeLists.txt\"", Str: "CMakeLists.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"README.md\"", Str: "README.md", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include/helper.hpp\"", Str: "include/helper.hpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/helper.cpp\"", Str: "src/helper.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/main.cpp\"", Str: "src/main.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test/TestHelper.cpp\"", Str: "test/TestHelper.cpp", Num: 0, Index: 0},
		},
		{
			{Type: gjson.Type(3), Raw: "\".vscode\"", Str: ".vscode", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include\"", Str: "include", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src\"", Str: "src", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test\"", Str: "test", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/c_cpp_properties.json\"", Str: ".vscode/c_cpp_properties.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/launch.json\"", Str: ".vscode/launch.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/tasks.json\"", Str: ".vscode/tasks.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"CMakeLists.txt\"", Str: "CMakeLists.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"README.md\"", Str: "README.md", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include/helper.hpp\"", Str: "include/helper.hpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"notes.txt\"", Str: "notes.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/helper.cpp\"", Str: "src/helper.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/main.cpp\"", Str: "src/main.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test/TestHelper.cpp\"", Str: "test/TestHelper.cpp", Num: 0, Index: 0},
		},
	}

	var goodRepoTypesResults = [][]gjson.Result{
		{
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
		},
		{
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
		},
	}

	// These should work
	for i, ghRepo := range goodRepoObjects {
		paths, types, err := ghRepo.GetTree(nil)
		assert.NoError(t, err)
		assert.NotNil(t, paths)
		assert.NotNil(t, types)
		assert.Equal(t, goodRepoPathsResults[i], paths)
		assert.Equal(t, goodRepoTypesResults[i], types)
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
