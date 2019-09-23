package project

import (
	"fmt"
	"os"
	"time"

	// Import sqlite3 driver (see func (setup *Setup) Run() error)
	"github.com/mattn/go-sqlite3"
)

// Project struct represents a project that will be build.
// Containing information about project name and label.
// The setup struct includes information about config paths and a open database connection.
type Project struct {
	ID         int
	Name       string
	InstallDir string
	Data       *Setup
}

// create starts the creation of a project.
// Returns an error on failure. Returns nil on success.
func (project *Project) new(projectID int) error {
	// Set id and installDir
	project.ID = projectID
	project.InstallDir = project.Data.Owd + project.Name

	// Create the project folder
	fmt.Println("> Creating project folder...")
	if err := project.createProjectFolder(); err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	if err := os.Chdir(project.Name); err != nil {
		return err
	}
	defer os.Chdir(project.Data.Owd)

	// Create subfolders
	fmt.Println("> Creating subfolders...")
	if err := project.createSubFolders(); err != nil {
		return err
	}

	// Create files
	fmt.Println("> Creating files...")
	if err := project.createFiles(); err != nil {
		return err
	}

	// Copy templates
	fmt.Println("> Copying templates...")
	if err := project.copyTemplates(); err != nil {
		return err
	}

	// Run scripts
	fmt.Println("> Running scripts...")
	return project.runScripts()
}

// track tracks a created project in the database
func (project *Project) track() error {
	t := time.Now().Local()
	_, err := project.Data.db.Exec(
		"INSERT INTO project(name, class_id, install_path, install_date, project_status_id) VALUES(?, ?, ?, ?, ?)",
		project.Name,
		project.ID,
		project.InstallDir,
		t,
		1,
	)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("project already exists")
		}
	}

	return err
}
