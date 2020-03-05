package item

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nikoksr/proji/pkg/config"

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
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	configPath, err := config.GetBaseConfigPath()
	if err != nil {
		t.Error(err)
	}

	tmpDir := os.TempDir()
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Error(err)
	}
	defer os.Chdir(originalCwd)

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
						{Destination: "README.md", Template: ""},
						{Destination: "exampleFolder/test.txt", Template: ""},
					},
					Scripts: make([]*Script, 0),
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
		err = test.proj.Create(tmpDir, configPath)
		assert.IsType(t, test.err, err)

		// Project folder should exist
		assert.DirExists(t, test.proj.Name)

		// Subfolders should exist
		for _, folder := range test.proj.Class.Folders {
			assert.DirExists(t, filepath.Join(test.proj.Name, folder.Destination))
		}

		// Project files should exist
		for _, file := range test.proj.Class.Files {
			assert.FileExists(t, filepath.Join(test.proj.Name, file.Destination))
		}

		// Compare old cwd to current cwd. Should be equal
		currentCwd, err := os.Getwd()
		if err != nil {
			t.FailNow()
		}

		// Special case for darwin systems.
		// On darwin systems /tmp is a symlink to /private/tmp which gets resolved by os.Getwd().
		// So the final path resulting from os.Getwd() differs from the original working directory
		// by the prefixed /private.
		if runtime.GOOS == "darwin" {
			assert.True(t, filepath.Join("/private", tmpDir) == currentCwd)
		} else {
			assert.True(t, tmpDir == currentCwd)
		}

		_ = os.RemoveAll(filepath.Join(tmpDir, test.proj.Name))
	}
}
