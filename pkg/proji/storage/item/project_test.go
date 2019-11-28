package item

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	class := NewClass("testclass", "tc", false)
	status := NewStatus(9999, "test", "This is a test status.", false)

	projExp := &Project{
		ID:          99,
		Name:        "test",
		InstallPath: "./local/",
		Class:       class,
		Status:      status,
	}

	projAct := NewProject(99, "test", "./local/", class, status)
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
		proj       *Project
		err        error
	}{
		{
			proj: &Project{
				Name:        "example",
				InstallPath: "",
				Class: &Class{
					Name:  "example",
					Label: "exp",
					Folders: []*Folder{
						{Destination: "exampleFolder/", Template: ""},
						{Destination: "foo/bar/", Template: ""},
					},
					Files: []*File{
						{Destination: "README.md", Template: "README.md"},
						{Destination: "exampleFolder/test.txt", Template: ""},
					},
					Scripts: []*Script{},
				},
				Status: &Status{
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

		_ = os.RemoveAll(originalCwd + "/" + test.proj.Name)
	}
}
