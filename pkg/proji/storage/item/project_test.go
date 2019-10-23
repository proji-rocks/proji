package item_test

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	class := item.NewClass("testclass", "tc")
	status := item.NewStatus(9999, "test", "This is a test status.")

	projExp := &item.Project{
		ID:          99,
		Name:        "test",
		InstallPath: "./local/",
		Class:       class,
		Status:      status,
	}

	projAct := item.NewProject(99, "test", "./local/", class, status)
	assert.Equal(t, projExp, projAct)
}

func TestProjectCreate(t *testing.T) {
	homeDir, err := homedir.Dir()
	if err != nil {
		t.FailNow()
	}
	configPath := homeDir + "/.config/proji/"

	tests := []struct {
		cwd        string
		configPath string
		proj       *item.Project
		err        error
	}{
		{
			proj: &item.Project{
				Name:        "example",
				InstallPath: "",
				Class: &item.Class{
					Name:  "example",
					Label: "exp",
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
				Status: &item.Status{
					ID:      1,
					Title:   "active",
					Comment: "This project is active.",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		originalCwd, err := os.Getwd()
		if err != nil {
			t.FailNow()
		}

		err = test.proj.Create(originalCwd, configPath)
		defer os.RemoveAll(originalCwd + "/" + test.proj.Name)
		assert.IsType(t, test.err, err)

		// Project folder should exist
		assert.DirExists(t, test.proj.Name)

		// Subfolders should exist
		for _, folder := range test.proj.Class.Folders {
			assert.DirExists(t, test.proj.Name+"/"+folder.Destination)
		}

		// Project files should exist
		for _, file := range test.proj.Class.Files {
			assert.FileExists(t, test.proj.Name+"/"+file.Destination)
		}

		// Compare old cwd to current cwd. Should be equal
		currentCwd, err := os.Getwd()
		if err != nil {
			t.FailNow()
		}
		assert.True(t, originalCwd == currentCwd)
	}
}
