package github

import (
	"testing"

	"github.com/tidwall/gjson"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/repo"

	"github.com/stretchr/testify/assert"
)

const ghAPIBase = "https://api.github.com/repos/"

var goodURLs = []string{
	"github.com/nikoksr/proji_test",
	"github.com/nikoksr/proji_test/tree/develop",
}

var goodRepoObjects = []repo.Importer{
	&github{
		apiBaseURI: ghAPIBase,
		userName:   "nikoksr",
		repoName:   "proji_test",
		branchName: "master",
		repoSHA:    "b4fc28f09ac57e314d27e9b9133b1ebc03bec2f1",
	},
	&github{
		apiBaseURI: ghAPIBase,
		userName:   "nikoksr",
		repoName:   "proji_test",
		branchName: "develop",
		repoSHA:    "f07d0b57cd6b468b331be03699f15faf4f9dd910",
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
	var badURLs = []string{
		"github.com/nikoksr/does-not-exist",
		"https://gitlab.com/nikoksr/does-not-exist",
		"https://google.com",
	}

	for _, URL := range badURLs {
		_, err := New(URL)
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
			{Type: gjson.Type(3), Raw: "\".vscode/c_cpp_properties.json\"", Str: ".vscode/c_cpp_properties.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/launch.json\"", Str: ".vscode/launch.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/tasks.json\"", Str: ".vscode/tasks.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"CMakeLists.txt\"", Str: "CMakeLists.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"README.md\"", Str: "README.md", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include\"", Str: "include", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include/helper.hpp\"", Str: "include/helper.hpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src\"", Str: "src", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/helper.cpp\"", Str: "src/helper.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/main.cpp\"", Str: "src/main.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test\"", Str: "test", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test/testHelper.cpp\"", Str: "test/testHelper.cpp", Num: 0, Index: 0},
		},
		{
			{Type: gjson.Type(3), Raw: "\".vscode\"", Str: ".vscode", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/c_cpp_properties.json\"", Str: ".vscode/c_cpp_properties.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/launch.json\"", Str: ".vscode/launch.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\".vscode/tasks.json\"", Str: ".vscode/tasks.json", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"CMakeLists.txt\"", Str: "CMakeLists.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"README.md\"", Str: "README.md", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include\"", Str: "include", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"include/helper.hpp\"", Str: "include/helper.hpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"notes.txt\"", Str: "notes.txt", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src\"", Str: "src", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/helper.cpp\"", Str: "src/helper.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"src/main.cpp\"", Str: "src/main.cpp", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test\"", Str: "test", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"test/testHelper.cpp\"", Str: "test/testHelper.cpp", Num: 0, Index: 0},
		},
	}

	var goodRepoTypesResults = [][]gjson.Result{
		{
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
		},
		{
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"tree\"", Str: "tree", Num: 0, Index: 0},
			{Type: gjson.Type(3), Raw: "\"blob\"", Str: "blob", Num: 0, Index: 0},
		},
	}

	// These should work
	for i, ghRepo := range goodRepoObjects {
		paths, types, err := ghRepo.GetTreePathsAndTypes()
		assert.NoError(t, err)
		assert.NotNil(t, paths)
		assert.NotNil(t, types)
		assert.Equal(t, goodRepoPathsResults[i], paths)
		assert.Equal(t, goodRepoTypesResults[i], types)
	}

	// These should fail
	var badRepoObjects = []repo.Importer{
		&github{
			apiBaseURI: "",
			userName:   "",
			repoName:   "",
			branchName: "",
			repoSHA:    "",
		},
		&github{
			apiBaseURI: ghAPIBase,
			userName:   "nikoksr",
			repoName:   "proji_test",
			branchName: "does_not_exist",
			repoSHA:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
	}

	for _, ghRepo := range badRepoObjects {
		paths, types, err := ghRepo.GetTreePathsAndTypes()
		assert.Error(t, err)
		assert.Nil(t, paths)
		assert.Nil(t, types)
	}
}

func Test_setRepoSHA(t *testing.T) {
	tests := []struct {
		name    string
		g       *github
		wantErr bool
		want    string
	}{
		{
			name: "Valid repo SHA",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji_test",
				branchName: "master",
				repoSHA:    "",
			},
			wantErr: false,
			want:    "b4fc28f09ac57e314d27e9b9133b1ebc03bec2f1",
		},
		{
			name: "Invalid repo",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji_private",
				branchName: "master",
				repoSHA:    "",
			},
			wantErr: true,
			want:    "",
		},
	}

	for _, test := range tests {
		err := test.g.setRepoSHA()
		assert.Equal(t, test.wantErr, err != nil)
		assert.Equal(t, test.want, test.g.repoSHA)
	}
}

func TestGetBranchName(t *testing.T) {
	tests := []struct {
		name string
		g    *github
		want string
	}{
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji",
				branchName: "master",
				repoSHA:    "",
			},
			want: "master",
		},
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji",
				branchName: "develop",
				repoSHA:    "",
			},
			want: "develop",
		},
	}
	for _, test := range tests {
		branch := test.g.GetBranchName()
		assert.Equal(t, test.want, branch, "%s\n", test.name)
	}
}

func TestGetRepoName(t *testing.T) {
	tests := []struct {
		name string
		g    *github
		want string
	}{
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji",
				branchName: "master",
				repoSHA:    "",
			},
			want: "proji",
		},
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "prinfo",
				branchName: "develop",
				repoSHA:    "",
			},
			want: "prinfo",
		},
	}
	for _, test := range tests {
		repo := test.g.GetRepoName()
		assert.Equal(t, test.want, repo, "%s\n", test.name)
	}
}

func TestGetUserName(t *testing.T) {
	tests := []struct {
		name string
		g    *github
		want string
	}{
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "nikoksr",
				repoName:   "proji",
				branchName: "master",
				repoSHA:    "",
			},
			want: "nikoksr",
		},
		{
			name: "",
			g: &github{
				apiBaseURI: ghAPIBase,
				userName:   "golang",
				repoName:   "go",
				branchName: "master",
				repoSHA:    "",
			},
			want: "golang",
		},
	}
	for _, test := range tests {
		user := test.g.GetUserName()
		assert.Equal(t, test.want, user, "%s\n", test.name)
	}
}
