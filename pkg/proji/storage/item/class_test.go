package item_test

import (
	"errors"
	"os"
	"testing"

	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/stretchr/testify/assert"
)

func TestNewClass(t *testing.T) {
	classExp := &item.Class{
		Name:      "test",
		Label:     "tst",
		IsDefault: false,
		Folders:   []*item.Folder{},
		Files:     []*item.File{},
		Scripts:   []*item.Script{},
	}

	classAct := item.NewClass("test", "tst", false)
	assert.Equal(t, classExp, classAct)
}

func TestClassImportFromConfig(t *testing.T) {
	tests := []struct {
		configName string
		class      *item.Class
		err        error
	}{
		{
			configName: "../../../../assets/examples/example-class-export.toml", class: &item.Class{
				Name:      "my-example",
				Label:     "mex",
				IsDefault: false,
				Folders: []*item.Folder{
					&item.Folder{Destination: "src/", Template: ""},
					&item.Folder{Destination: "docs/", Template: ""},
					&item.Folder{Destination: "tests/", Template: ""},
				},
				Files: []*item.File{
					&item.File{Destination: "src/main.py", Template: ""},
					&item.File{Destination: "README.md", Template: ""},
				},
				Scripts: []*item.Script{
					&item.Script{
						Name:       "init_virtualenv.sh",
						Type:       "post",
						ExecNumber: 1,
						RunAsSudo:  false,
						Args:       []string{},
					},
					&item.Script{
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
			class: &item.Class{
				Name:      "",
				Label:     "",
				IsDefault: false,
				Folders:   []*item.Folder{},
				Files:     []*item.File{},
				Scripts:   []*item.Script{},
			},
			err: errors.New(""),
		},
	}

	for _, test := range tests {
		c := item.NewClass("", "", false)
		err := c.ImportFromConfig(test.configName)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.class, c)
	}
}

func TestClassImportFromDirectory(t *testing.T) {
	tests := []struct {
		baseName string
		folders  []*item.Folder
		files    []*item.File
		class    *item.Class
		err      error
	}{
		{
			baseName: "new-project",
			folders: []*item.Folder{
				&item.Folder{Destination: "new-project", Template: ""},
				&item.Folder{Destination: "new-project/test", Template: ""},
				&item.Folder{Destination: "new-project/cmd", Template: ""},
				&item.Folder{Destination: "new-project/cmd/base", Template: ""},
				&item.Folder{Destination: "new-project/docs", Template: ""},
			},
			files: []*item.File{
				&item.File{Destination: "new-project/test.txt", Template: ""},
				&item.File{Destination: "new-project/README.md", Template: ""},
				&item.File{Destination: "new-project/cmd/main.go", Template: ""},
				&item.File{Destination: "new-project/test/main_test.go", Template: ""},
			},
			class: &item.Class{
				Name:      "new-project",
				Label:     "np",
				IsDefault: false,
				Folders: []*item.Folder{
					&item.Folder{Destination: "cmd", Template: ""},
					&item.Folder{Destination: "cmd/base", Template: ""},
					&item.Folder{Destination: "docs", Template: ""},
					&item.Folder{Destination: "test", Template: ""},
				},
				Files: []*item.File{
					&item.File{Destination: "README.md", Template: ""},
					&item.File{Destination: "cmd/main.go", Template: ""},
					&item.File{Destination: "test/main_test.go", Template: ""},
					&item.File{Destination: "test.txt", Template: ""},
				},
				Scripts: []*item.Script{},
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
		c := item.NewClass("", "", false)
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
		class      *item.Class
		configName string
		err        error
	}{
		{
			class: &item.Class{
				Name:      "example",
				Label:     "exp",
				IsDefault: false,
				Folders: []*item.Folder{
					&item.Folder{Destination: "exampleFolder/", Template: ""},
					&item.Folder{Destination: "foo/bar/", Template: ""},
				},
				Files: []*item.File{
					&item.File{Destination: "README.md", Template: "README.md"},
					&item.File{Destination: "exampleFolder/test.txt", Template: ""},
				},
				Scripts: []*item.Script{},
			},
			configName: "./proji-example.toml",
			err:        nil,
		},
	}

	for _, test := range tests {
		configName, err := test.class.Export(".")
		defer os.Remove(configName)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.configName, configName)
		assert.FileExists(t, test.configName, "Cannot find the exported config file.")
	}
}
