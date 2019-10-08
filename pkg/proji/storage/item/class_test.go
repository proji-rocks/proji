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
		Name:    "test",
		Label:   "tst",
		Folders: map[string]string{},
		Files:   map[string]string{},
		Scripts: map[string]bool{},
	}

	classAct := item.NewClass("test", "tst")
	assert.Equal(t, classExp, classAct)
}

func TestClassImportData(t *testing.T) {
	tests := []struct {
		configName string
		class      *item.Class
		err        error
	}{
		{
			configName: "../../../../configs/example-class-export.toml", class: &item.Class{
				Name:  "example",
				Label: "exp",
				Folders: map[string]string{
					"exampleFolder/": "",
					"foo/bar/":       "",
				},
				Files: map[string]string{
					"README.md":              "README.md",
					"exampleFolder/test.txt": "",
				},
				Scripts: map[string]bool{},
			},
			err: nil,
		},
		{
			configName: "example.yaml",
			class: &item.Class{
				Name:    "",
				Label:   "",
				Folders: map[string]string{},
				Files:   map[string]string{},
				Scripts: map[string]bool{},
			},
			err: errors.New(""),
		},
	}

	for _, test := range tests {
		c := item.NewClass("", "")
		err := c.ImportData(test.configName)
		assert.IsType(t, test.err, err)
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
				Name:  "example",
				Label: "exp",
				Folders: map[string]string{
					"exampleFolder/": "",
					"foo/bar/":       "",
				},
				Files: map[string]string{
					"README.md":              "README.md",
					"exampleFolder/test.txt": "",
				},
				Scripts: map[string]bool{},
			},
			configName: "proji-example.toml",
			err:        nil,
		},
	}

	for _, test := range tests {
		configName, err := test.class.Export()
		defer os.Remove(configName)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.configName, configName)
		assert.FileExists(t, test.configName, "Cannot find the exported config file.")
	}
}
