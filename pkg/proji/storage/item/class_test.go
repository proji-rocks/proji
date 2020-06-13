package item

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/nikoksr/proji/pkg/config"

	"github.com/nikoksr/proji/pkg/proji/repo/github"

	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/nikoksr/proji/pkg/proji/repo/gitlab"

	"github.com/nikoksr/proji/pkg/helper"

	"github.com/stretchr/testify/assert"
)

var goodURLs = []*url.URL{
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test"},
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/tree/develop"},
	{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test"},
	{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test/tree/develop"},
}

var badURLs = []*url.URL{
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/does-not-exist"},
	{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/tree/dead-branch"},
	{Scheme: "https", Host: "google.com", Path: ""},
}

var apiAuth = config.APIAuthentication{
	GHToken: os.Getenv("PROJI_AUTH_GH_TOKEN"),
	GLToken: os.Getenv("PROJI_AUTH_GL_TOKEN"),
}

func TestNewClass(t *testing.T) {
	classExp := &Class{
		Name:      "test",
		Label:     "tst",
		IsDefault: false,
		Folders:   make([]*Folder, 0),
		Files:     make([]*File, 0),
		Scripts:   make([]*Script, 0),
	}

	classAct := NewClass("test", "tst", false)
	assert.Equal(t, classExp, classAct)
}

func TestClass_ImportConfig(t *testing.T) {
	tests := []struct {
		configName string
		class      *Class
		err        error
	}{
		{
			configName: "../../../../assets/examples/example-class-export.toml",
			class: &Class{
				Name:      "my-example",
				Label:     "mex",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: "src/", Template: ""},
					{Destination: "docs/", Template: ""},
					{Destination: "tests/", Template: ""},
				},
				Files: []*File{
					{Destination: "src/main.py", Template: ""},
					{Destination: "README.md", Template: ""},
				},
				Scripts: []*Script{
					{
						Name:       "init_virtualenv.sh",
						Type:       "post",
						ExecNumber: 1,
						RunAsSudo:  false,
						Args:       make([]string, 0),
					},
					{
						Name:       "init_git.sh",
						Type:       "post",
						ExecNumber: 2,
						RunAsSudo:  false,
						Args:       make([]string, 0),
					},
				},
			},
			err: nil,
		},
		{
			configName: "example.yaml",
			class: &Class{
				Name:      "",
				Label:     "",
				IsDefault: false,
				Folders:   make([]*Folder, 0),
				Files:     make([]*File, 0),
				Scripts:   make([]*Script, 0),
			},
			err: errors.New(""),
		},
	}

	for _, test := range tests {
		c := NewClass("", "", false)
		err := c.ImportConfig(test.configName)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.class, c)
	}
}

func TestClass_ImportFolderStructure(t *testing.T) {
	tmpDir := os.TempDir()

	tests := []struct {
		basePath string
		folders  []*Folder
		files    []*File
		class    *Class
		err      error
	}{
		{
			basePath: filepath.Join(tmpDir, "/proji/new-project"),
			folders: []*Folder{
				{Destination: "test", Template: ""},
				{Destination: "cmd", Template: ""},
				{Destination: "cmd/base", Template: ""},
				{Destination: "docs", Template: ""},
			},
			files: []*File{
				{Destination: "test.txt", Template: ""},
				{Destination: "README.md", Template: ""},
				{Destination: "cmd/main.go", Template: ""},
				{Destination: "test/main_test.go", Template: ""},
			},
			class: &Class{
				Name:      "new-project",
				Label:     "np",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: "cmd", Template: ""},
					{Destination: "cmd/base", Template: ""},
					{Destination: "docs", Template: ""},
					{Destination: "test", Template: ""},
				},
				Files: []*File{
					{Destination: "README.md", Template: ""},
					{Destination: "cmd/main.go", Template: ""},
					{Destination: "test/main_test.go", Template: ""},
					{Destination: "test.txt", Template: ""},
				},
				Scripts: make([]*Script, 0),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		for _, dir := range test.folders {
			assert.NoError(t, os.MkdirAll(filepath.Join(test.basePath, dir.Destination), os.ModePerm))
		}
		for _, file := range test.files {
			_, err := os.Create(filepath.Join(test.basePath, file.Destination))
			assert.NoError(t, err)
		}

		c := NewClass("", "", false)
		assert.NoError(t, c.ImportFolderStructure(test.basePath, make([]string, 0)))
		conf, err := c.Export(tmpDir)
		assert.NoError(t, err)
		assert.NoError(t, c.ImportConfig(conf))
		assert.Equal(t, test.class, c)

		// Clean up
		_ = os.Remove(conf)
		_ = os.RemoveAll(test.basePath)
	}
}

func TestClass_ImportRepoStructure(t *testing.T) {
	helper.SkipNetworkBasedTests(t)

	tests := []struct {
		URL   *url.URL
		class *Class
		err   error
	}{
		{
			URL: &url.URL{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test"},
			class: &Class{
				Name:      "proji-test",
				Label:     "pt",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: ".vscode", Template: ""},
					{Destination: "include", Template: ""},
					{Destination: "src", Template: ""},
					{Destination: "test", Template: ""},
				},
				Files: []*File{
					{Destination: ".vscode/c_cpp_properties.json", Template: ""},
					{Destination: ".vscode/launch.json", Template: ""},
					{Destination: ".vscode/tasks.json", Template: ""},
					{Destination: "CMakeLists.txt", Template: ""},
					{Destination: "README.md", Template: ""},
					{Destination: "include/helper.hpp", Template: ""},
					{Destination: "src/helper.cpp", Template: ""},
					{Destination: "src/main.cpp", Template: ""},
					{Destination: "test/testHelper.cpp", Template: ""},
				},
				Scripts: make([]*Script, 0),
			},
			err: nil,
		},
		{
			URL: &url.URL{Scheme: "https", Host: "github.com", Path: "/nikoksr/proji-test/tree/develop"},
			class: &Class{
				Name:      "proji-test",
				Label:     "pt",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: ".vscode", Template: ""},
					{Destination: "include", Template: ""},
					{Destination: "src", Template: ""},
					{Destination: "test", Template: ""},
				},
				Files: []*File{
					{Destination: ".vscode/c_cpp_properties.json", Template: ""},
					{Destination: ".vscode/launch.json", Template: ""},
					{Destination: ".vscode/tasks.json", Template: ""},
					{Destination: "CMakeLists.txt", Template: ""},
					{Destination: "README.md", Template: ""},
					{Destination: "include/helper.hpp", Template: ""},
					{Destination: "notes.txt", Template: ""},
					{Destination: "src/helper.cpp", Template: ""},
					{Destination: "src/main.cpp", Template: ""},
					{Destination: "test/testHelper.cpp", Template: ""},
				},
				Scripts: make([]*Script, 0),
			},
			err: nil,
		},
		{
			URL: &url.URL{Scheme: "https", Host: "gitlab.com", Path: "/nikoksr/proji-test"},
			class: &Class{
				Name:      "proji-test",
				Label:     "pt",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: ".vscode", Template: ""},
					{Destination: "include", Template: ""},
					{Destination: "src", Template: ""},
					{Destination: "test", Template: ""},
				},
				Files: []*File{
					{Destination: ".vscode/c_cpp_properties.json", Template: ""},
					{Destination: ".vscode/launch.json", Template: ""},
					{Destination: ".vscode/tasks.json", Template: ""},
					{Destination: "CMakeLists.txt", Template: ""},
					{Destination: "README.md", Template: ""},
					{Destination: "include/helper.hpp", Template: ""},
					{Destination: "src/helper.cpp", Template: ""},
					{Destination: "src/main.cpp", Template: ""},
					{Destination: "test/TestHelper.cpp", Template: ""},
				},
				Scripts: make([]*Script, 0),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		c := NewClass("", "", false)
		importer, err := GetRepoImporterFromURL(test.URL, &apiAuth)
		if err != nil {
			t.Errorf("failed getting repo importer for URL %s", test.URL.String())
		}
		assert.NoError(t, c.ImportRepoStructure(importer, nil))
		assert.Equal(t, test.class, c)
	}
}

func TestClass_Export(t *testing.T) {
	tmpDir := os.TempDir()

	tests := []struct {
		class      *Class
		configPath string
		err        error
	}{
		{
			class: &Class{
				Name:      "example",
				Label:     "exp",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: "exampleFolder/", Template: ""},
					{Destination: "foo/bar/", Template: ""},
				},
				Files: []*File{
					{Destination: "README.md", Template: "README.md"},
					{Destination: "exampleFolder/test.txt", Template: ""},
				},
				Scripts: make([]*Script, 0),
			},
			configPath: filepath.Join(tmpDir, "/proji-example.toml"),
			err:        nil,
		},
	}

	for _, test := range tests {
		configPath, err := test.class.Export(tmpDir)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.configPath, configPath)
		assert.FileExists(t, configPath, "Cannot find the exported config file.")
		_ = os.Remove(configPath)
	}
}

func TestClass_isEmpty(t *testing.T) {
	tests := []struct {
		class   *Class
		isEmpty bool
	}{
		{
			class: &Class{
				Name:      "my-example",
				Label:     "mex",
				IsDefault: false,
				Folders: []*Folder{
					{Destination: "src/", Template: ""},
					{Destination: "docs/", Template: ""},
					{Destination: "tests/", Template: ""},
				},
				Files: []*File{
					{Destination: "src/main.py", Template: ""},
					{Destination: "README.md", Template: ""},
				},
				Scripts: []*Script{
					{
						Name:       "init_virtualenv.sh",
						Type:       "post",
						ExecNumber: 1,
						RunAsSudo:  false,
						Args:       make([]string, 0),
					},
					{
						Name:       "init_git.sh",
						Type:       "post",
						ExecNumber: 2,
						RunAsSudo:  false,
						Args:       make([]string, 0),
					},
				},
			},
			isEmpty: false,
		},
		{
			class: &Class{
				Name:      "blabla",
				Label:     "bl",
				IsDefault: false,
				Folders:   make([]*Folder, 0),
				Files:     make([]*File, 0),
				Scripts:   make([]*Script, 0),
			},
			isEmpty: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.class.isEmpty(), test.isEmpty)
	}
}

func TestGetRepoImporterFromURL(t *testing.T) {
	type args struct {
		URL *url.URL
	}

	tests := []struct {
		name    string
		args    args
		want    repo.Importer
		wantErr bool
	}{
		{
			name: "Test Importer 1",
			args: args{URL: goodURLs[0]},
			want: &github.GitHub{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test",
				BranchName: "master",
			},
			wantErr: false,
		},
		{
			name: "Test Importer 2",
			args: args{URL: goodURLs[1]},
			want: &github.GitHub{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test",
				BranchName: "develop",
			},
			wantErr: false,
		},
		{
			name: "Test Importer 3",
			args: args{URL: goodURLs[2]},
			want: &gitlab.GitLab{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test",
				BranchName: "master",
			},
			wantErr: false,
		},
		{
			name: "Test Importer 4",
			args: args{URL: goodURLs[3]},
			want: &gitlab.GitLab{
				OwnerName:  "nikoksr",
				RepoName:   "proji-test",
				BranchName: "develop",
			},
			wantErr: false,
		},
		{
			name:    "Test Importer 5",
			args:    args{URL: badURLs[0]},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test Importer 6",
			args:    args{URL: badURLs[1]},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test Importer 7",
			args:    args{URL: badURLs[2]},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRepoImporterFromURL(tt.args.URL, &apiAuth)
			assert.Equal(t, err != nil, tt.wantErr, "GetRepoImporterFromURL() error = %v, wantErr %v", err, tt.wantErr)

			if got != nil && tt.want != nil {
				assert.Equal(t, tt.want.Owner(), got.Owner(), tt.name)
				assert.Equal(t, tt.want.Repo(), got.Repo(), tt.name)
				assert.Equal(t, tt.want.Branch(), got.Branch(), tt.name)
			}
		})
	}
}

func Test_pickLabel(t *testing.T) {
	type args struct {
		className string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Pick Label 1",
			args: args{className: "myTestClass"},
			want: "mtc",
		},
		{
			name: "Test Pick Label 2",
			args: args{className: "my-test-class"},
			want: "mtc",
		},
		{
			name: "Test Pick Label 3",
			args: args{className: "my.test.class"},
			want: "mtc",
		},
		{
			name: "Test Pick Label 4",
			args: args{className: "my_test_class"},
			want: "mtc",
		},
		{
			name: "Test Pick Label 5",
			args: args{className: "my%20test%20class"},
			want: "mtc",
		},
		{
			name: "Test Pick Label 6",
			args: args{className: "mytestclass"},
			want: "mts",
		},
		{
			name: "Test Pick Label 7",
			args: args{className: "sjcsdhfklhaslcsdflsancshdkljfjalksfjnsvnslkd"},
			want: "ssd",
		},
		{
			name: "Test Pick Label 8",
			args: args{className: "s"},
			want: "s",
		},
		{
			name: "Test Pick Label 9",
			args: args{className: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickLabel(tt.args.className)
			assert.Equal(t, tt.want, got, "pickLabel() = %v, want %v", got, tt.want)
		})
	}
}
