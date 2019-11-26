package item

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClass(t *testing.T) {
	classExp := &Class{
		Name:      "test",
		Label:     "tst",
		IsDefault: false,
		Folders:   []*Folder{},
		Files:     []*File{},
		Scripts:   []*Script{},
	}

	classAct := NewClass("test", "tst", false)
	assert.Equal(t, classExp, classAct)
}

func TestClassImportFromConfig(t *testing.T) {
	tests := []struct {
		configName string
		class      *Class
		err        error
	}{
		{
			configName: "../../../../assets/examples/example-class-export.toml", class: &Class{
				Name:      "my-example",
				Label:     "mex",
				IsDefault: false,
				Folders: []*Folder{
					&Folder{Destination: "src/", Template: ""},
					&Folder{Destination: "docs/", Template: ""},
					&Folder{Destination: "tests/", Template: ""},
				},
				Files: []*File{
					&File{Destination: "src/main.py", Template: ""},
					&File{Destination: "README.md", Template: ""},
				},
				Scripts: []*Script{
					&Script{
						Name:       "init_virtualenv.sh",
						Type:       "post",
						ExecNumber: 1,
						RunAsSudo:  false,
						Args:       []string{},
					},
					&Script{
						Name:       "init_git.sh",
						Type:       "post",
						ExecNumber: 2,
						RunAsSudo:  false,
						Args:       []string{},
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
				Folders:   []*Folder{},
				Files:     []*File{},
				Scripts:   []*Script{},
			},
			err: errors.New(""),
		},
	}

	for _, test := range tests {
		c := NewClass("", "", false)
		err := c.ImportFromConfig(test.configName)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.class, c)
	}
}

func TestClassImportFromDirectory(t *testing.T) {
	tests := []struct {
		baseName string
		folders  []*Folder
		files    []*File
		class    *Class
		err      error
	}{
		{
			baseName: "new-project",
			folders: []*Folder{
				&Folder{Destination: "new-project", Template: ""},
				&Folder{Destination: "new-project/test", Template: ""},
				&Folder{Destination: "new-project/cmd", Template: ""},
				&Folder{Destination: "new-project/cmd/base", Template: ""},
				&Folder{Destination: "new-project/docs", Template: ""},
			},
			files: []*File{
				&File{Destination: "new-project/test.txt", Template: ""},
				&File{Destination: "new-project/README.md", Template: ""},
				&File{Destination: "new-project/cmd/main.go", Template: ""},
				&File{Destination: "new-project/test/main_test.go", Template: ""},
			},
			class: &Class{
				Name:      "new-project",
				Label:     "np",
				IsDefault: false,
				Folders: []*Folder{
					&Folder{Destination: "cmd", Template: ""},
					&Folder{Destination: "cmd/base", Template: ""},
					&Folder{Destination: "docs", Template: ""},
					&Folder{Destination: "test", Template: ""},
				},
				Files: []*File{
					&File{Destination: "README.md", Template: ""},
					&File{Destination: "cmd/main.go", Template: ""},
					&File{Destination: "test/main_test.go", Template: ""},
					&File{Destination: "test.txt", Template: ""},
				},
				Scripts: []*Script{},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		for _, dir := range test.folders {
			assert.NoError(t, os.Mkdir(dir.Destination, 0755))
		}
		for _, file := range test.files {
			_, err := os.Create(file.Destination)
			assert.NoError(t, err)
		}
		defer os.RemoveAll(test.baseName)
		c := NewClass("", "", false)
		assert.NoError(t, c.ImportFromDirectory(test.baseName, []string{}))
		conf, err := c.Export(".")
		defer os.Remove(conf)
		assert.NoError(t, err)
		assert.NoError(t, c.ImportFromConfig(conf))
		assert.Equal(t, test.class, c)
	}
}

func TestClassExport(t *testing.T) {
	tests := []struct {
		class      *Class
		configName string
		err        error
	}{
		{
			class: &Class{
				Name:      "example",
				Label:     "exp",
				IsDefault: false,
				Folders: []*Folder{
					&Folder{Destination: "exampleFolder/", Template: ""},
					&Folder{Destination: "foo/bar/", Template: ""},
				},
				Files: []*File{
					&File{Destination: "README.md", Template: "README.md"},
					&File{Destination: "exampleFolder/test.txt", Template: ""},
				},
				Scripts: []*Script{},
			},
			configName: "./proji-example.toml",
			err:        nil,
		},
	}

	for _, test := range tests {
		configName, err := test.class.Export(".")
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.configName, configName)
		assert.FileExists(t, test.configName, "Cannot find the exported config file.")
		os.Remove(configName)
	}
}

func TestClassIsEmpty(t *testing.T) {
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
					&Folder{Destination: "src/", Template: ""},
					&Folder{Destination: "docs/", Template: ""},
					&Folder{Destination: "tests/", Template: ""},
				},
				Files: []*File{
					&File{Destination: "src/main.py", Template: ""},
					&File{Destination: "README.md", Template: ""},
				},
				Scripts: []*Script{
					&Script{
						Name:       "init_virtualenv.sh",
						Type:       "post",
						ExecNumber: 1,
						RunAsSudo:  false,
						Args:       []string{},
					},
					&Script{
						Name:       "init_git.sh",
						Type:       "post",
						ExecNumber: 2,
						RunAsSudo:  false,
						Args:       []string{},
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
				Folders:   []*Folder{},
				Files:     []*File{},
				Scripts:   []*Script{},
			},
			isEmpty: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.class.isEmpty(), test.isEmpty)
	}
}
