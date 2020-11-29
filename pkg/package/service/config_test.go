package packageservice

import (
	"fmt"
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	packageAsToml = `name = "golang"
label = "go"
description = "Test package"

[[template]]
  is_file = false
  destination = "cmd"
  path = ""
  description = "cmd folder"

[[template]]
  is_file = true
  destination = "main.go"
  path = "main.go"
  description = "entry file"

[[plugin]]
  path = "git-init.lua"
  exec_number = 1
  description = ""
`

	packageAsJson = `{"name":"golang","label":"go","description":"Test package","template":[{"is_file":false,"destination":"cmd","path":"","description":"cmd folder"},{"is_file":true,"destination":"main.go","path":"main.go","description":"entry file"}],"plugin":[{"path":"git-init.lua","exec_number":1,"description":""}]}
`
)

func TestExportPackageToConfig(t *testing.T) {
	dest := t.TempDir()
	err := os.MkdirAll(dest, os.ModePerm)
	assert.NoError(t, err, "Unable to create temporary file")

	ps := packageService{
		authentication: nil,
		packageStore:   nil,
	}

	type args struct {
		json bool
		pkg  domain.Package
		dest string
	}
	tests := []struct {
		name             string
		expectedFilename string
		args             args
		expectedOutput   string
	}{
		{
			name:             "Create TOML export",
			args:             args{json: false, pkg: getPackage(), dest: dest},
			expectedFilename: "proji-golang.toml",
			expectedOutput:   packageAsToml,
		},
		{
			name:             "Create JSON export",
			args:             args{json: true, pkg: getPackage(), dest: dest},
			expectedFilename: "proji-golang.json",
			expectedOutput:   packageAsJson,
		},
	}

	for _, test := range tests {
		savedTo, err := ps.ExportPackageToConfig(test.args.pkg, test.args.dest, test.args.json)
		expectedLocation := fmt.Sprintf("%s/%s", test.args.dest, test.expectedFilename)
		assert.NoError(t, err, test.name)
		assert.Equal(t, expectedLocation, savedTo, test.name)

		file, err := os.Open(expectedLocation)
		assert.NoError(t, err, test.name)
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		assert.Equal(t, test.expectedOutput, string(bytes), test.name)
	}
}

func TestImportPackageFromConfig(t *testing.T) {
	dest := t.TempDir()
	tomlPath := filepath.Join(dest, "proji-golang.toml")
	err := ioutil.WriteFile(tomlPath, []byte(packageAsToml), os.ModePerm)
	assert.NoError(t, err, "error when creating temp file %s", tomlPath)
	jsonPath := filepath.Join(dest, "proji-golang.json")
	err = ioutil.WriteFile(jsonPath, []byte(packageAsJson), os.ModePerm)
	assert.NoError(t, err, "error when creating temp file %s", jsonPath)

	ps := packageService{
		authentication: nil,
		packageStore:   nil,
	}

	testFiles := []string{tomlPath, jsonPath}

	expectedPkg := getPackage()

	for _, file := range testFiles {
		pkg, err := ps.ImportPackageFromConfig(file)
		assert.NoError(t, err, "error when importing %s", file)
		assert.Equal(t, expectedPkg.Name, pkg.Name, "incorrect import from %s", file)
		assert.Equal(t, expectedPkg.Label, pkg.Label, "incorrect import from %s", file)
		assert.Equal(t, expectedPkg.Description, pkg.Description, "incorrect import from %s", file)
		for i, tmpl := range expectedPkg.Templates {
			assert.Equal(t, tmpl.Description, pkg.Templates[i].Description, "incorrect import from %s", file)
			assert.Equal(t, tmpl.IsFile, pkg.Templates[i].IsFile, "incorrect import from %s", file)
			assert.Equal(t, tmpl.Destination, pkg.Templates[i].Destination, "incorrect import from %s", file)
			assert.Equal(t, tmpl.Path, pkg.Templates[i].Path, "incorrect import from %s", file)
		}
		for i, plugin := range expectedPkg.Plugins {
			assert.Equal(t, plugin.Description, pkg.Plugins[i].Description, "incorrect import from %s", file)
			assert.Equal(t, plugin.ExecNumber, pkg.Plugins[i].ExecNumber, "incorrect import from %s", file)
			assert.Equal(t, plugin.Path, pkg.Plugins[i].Path, "incorrect import from %s", file)
		}
	}
}

func getPackage() domain.Package {
	tmpl1 := domain.Template{
		ID:          2,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsFile:      false,
		Destination: "cmd",
		Path:        "",
		Description: "cmd folder",
	}
	tmpl2 := domain.Template{
		ID:          3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsFile:      true,
		Destination: "main.go",
		Path:        "main.go",
		Description: "entry file",
	}
	plugin := domain.Plugin{
		ID:          4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Path:        "git-init.lua",
		ExecNumber:  1,
		Description: "",
	}

	templates := []*domain.Template{&tmpl1, &tmpl2}
	plugins := []*domain.Plugin{&plugin}

	return domain.Package{
		ID:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        "golang",
		Label:       "go",
		Description: "Test package",
		Templates:   templates,
		Plugins:     plugins,
	}
}
