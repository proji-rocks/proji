package project

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	// Import sqlite3 driver (see func (setup *Setup) Run() error)
	_ "github.com/mattn/go-sqlite3"
	"github.com/nikoksr/proji/internal/app/helper"

	"github.com/otiai10/copy"
)

// CreateProject will create projects.
// It will create directories and files, copy templates and run scripts.
func CreateProject(label string, projects []string) error {
	// TODO: Load values from a config file
	configDir := helper.GetConfigDir()
	databaseName, ok := viper.Get("database.name").(string)

	if ok != true {
		return errors.New("could not read database name from config file")
	}

	// Get current working directory
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	// Create setup
	label = strings.ToLower(label)
	newSetup := Setup{Owd: cwd, ConfigDir: configDir, DatabaseName: databaseName, Label: label}
	err = newSetup.init()
	if err != nil {
		return err
	}
	defer newSetup.stop()

	// Projects loop
	for _, projectName := range projects {
		fmt.Println(helper.ProjectHeader(projectName))
		newProject := Project{Name: projectName, Data: &newSetup}
		err = newProject.create()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

// Setup contains necessary informations for the creation of a project.
// Owd is the Origin Working Directory.
type Setup struct {
	Owd          string
	ConfigDir    string
	DatabaseName string
	Label        string
	templatesDir string
	scriptsDir   string
	dbDir        string
	db           *sql.DB
	projectID    string
}

// init initializes the setup struct. Creates a database connection and defines default directores.
func (setup *Setup) init() error {
	// Set dirs
	setup.dbDir = setup.ConfigDir + "db/"
	setup.templatesDir = setup.ConfigDir + "templates/"
	setup.scriptsDir = setup.ConfigDir + "scripts/"

	// Connect to database
	db, err := sql.Open("sqlite3", setup.dbDir+setup.DatabaseName)
	if err != nil {
		return err
	}
	setup.db = db

	// Check if label is supported
	err = setup.isLabelSupported()
	if err != nil {
		return err
	}
	return nil
}

// stop cleanly stops the running Setup instance.
// Currently it's only closing its open database connection.
func (setup *Setup) stop() {
	// Close database connection
	if setup.db != nil {
		setup.db.Close()
	}
}

// isLabelSupported checks if the given label is found in the database.
// Returns nil if found, returns error if not found
func (setup *Setup) isLabelSupported() error {
	stmt, err := setup.db.Prepare("SELECT project_class_id FROM class_label WHERE label = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(setup.Label).Scan(&id)
	if err != nil {
		return err
	}
	setup.projectID = id
	return nil
}

// Project struct represents a project that will be build.
// Containing information about project name and label.
// The setup struct includes information about config paths and a open database connection.
type Project struct {
	id   string
	Name string
	Data *Setup
}

// create starts the creation of a project.
// Returns an error on failure. Returns nil on success.
func (project *Project) create() error {
	// Create the project folder
	fmt.Println("> Creating project folder...")
	err := project.createProjectFolder()
	if err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	err = os.Chdir(project.Name)
	if err != nil {
		return err
	}
	defer os.Chdir(project.Data.Owd)

	// Create subfolders
	fmt.Println("> Creating subfolders...")
	err = project.createSubFolders()
	if err != nil {
		return err
	}

	// Create files
	fmt.Println("> Creating files...")
	err = project.createFiles()
	if err != nil {
		return err
	}

	// Copy templates
	fmt.Println("> Copying templates...")
	err = project.copyTemplates()
	if err != nil {
		return err
	}

	// Run scripts
	fmt.Println("> Running scripts...")
	err = project.runScripts()
	if err != nil {
		return err
	}
	return nil
}

// createProjectFolder tries to create the main project folder.
// Returns an error on failure.
func (project *Project) createProjectFolder() error {
	err := os.Mkdir(project.Name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// createSubFolders queries all subfolder from the database related to the projectId.
// Tries to create all of the subfolders in the projectfolder.
// Returns error on failure. Returns nil on success.
func (project *Project) createSubFolders() error {
	// Query subfolders
	stmt, err := project.Data.db.Prepare("SELECT target_path FROM class_folder WHERE (project_class_id = ? OR project_class_id IS NULL) AND template_name IS NULL")
	if err != nil {
		return err
	}
	defer stmt.Close()

	subFolders, err := stmt.Query(project.Data.projectID)
	if err != nil {
		return err
	}
	defer subFolders.Close()

	// Create subfolders
	for subFolders.Next() {
		var subFolder string
		err = subFolders.Scan(&subFolder)
		if err != nil {
			return err
		}
		err = os.MkdirAll(subFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// createFiles queries all files from the database related to the projectId.
// Tries to create all of the files in the projectfolder.
// Returns error on failure. Returns nil on success.
func (project *Project) createFiles() error {
	// Query files
	stmt, err := project.Data.db.Prepare("SELECT target_path FROM class_file WHERE (project_class_id = ? OR project_class_id IS NULL) AND template_name IS NULL")
	if err != nil {
		return err
	}
	defer stmt.Close()

	files, err := stmt.Query(project.Data.projectID)
	if err != nil {
		return err
	}
	defer files.Close()

	// Create files
	for files.Next() {
		var file string
		err = files.Scan(&file)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// copyTemplates queries all templates from the database related to the projectId.
// Tries to copy all of the templates into the projectfolder.
// Returns error on failure. Returns nil on success.
func (project *Project) copyTemplates() error {
	// Query template folders
	stmt, err := project.Data.db.Prepare(
		"SELECT target_path, template_name FROM class_folder WHERE (project_class_id = ? OR project_class_id IS NULL) AND template_name IS NOT NULL")
	if err != nil {
		return err
	}
	defer stmt.Close()

	folders, err := stmt.Query(project.Data.projectID)
	if err != nil {
		return err
	}
	defer folders.Close()

	// Copy template files
	for folders.Next() {
		var target, src string
		err = folders.Scan(&target, &src)
		if err != nil {
			return err
		}
		src = project.Data.templatesDir + src
		err := copy.Copy(src, target)
		if err != nil {
			return err
		}
	}

	// Query template files
	stmt, err = project.Data.db.Prepare(
		"SELECT target_path, template_name FROM class_file WHERE (project_class_id = ? OR project_class_id IS NULL) AND template_name IS NOT NULL")
	if err != nil {
		return err
	}

	files, err := stmt.Query(project.Data.projectID)
	if err != nil {
		return err
	}
	defer files.Close()

	// Copy template files
	for files.Next() {
		var target, src string
		err = files.Scan(&target, &src)
		if err != nil {
			return err
		}
		src = project.Data.templatesDir + src
		err := copy.Copy(src, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// runScripts queries all scripts from the database related to the projectId.
// Tries to execute all scripts.
// Returns error on failure. Returns nil on success.
func (project *Project) runScripts() error {
	// Query scripts
	stmt, err := project.Data.db.Prepare("SELECT script_name, run_as_sudo FROM class_script WHERE project_class_id is NULL OR project_class_id = ? ORDER BY project_class_id DESC")
	if err != nil {
		return err
	}
	defer stmt.Close()

	scripts, err := stmt.Query(project.Data.projectID)
	if err != nil {
		return err
	}
	defer scripts.Close()

	// Create scripts
	for scripts.Next() {
		var script string
		var runAsSudo int
		err = scripts.Scan(&script, &runAsSudo)
		if err != nil {
			return err
		}
		script = project.Data.scriptsDir + script
		err = exec.Command(script).Run()
		if err != nil {
			return err
		}
	}
	return nil
}
