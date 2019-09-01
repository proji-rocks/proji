package project

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	// Import sqlite3 driver (see func (setup *Setup) Run() error)
	"github.com/mattn/go-sqlite3"
	"github.com/nikoksr/proji/internal/app/helper"

	"github.com/otiai10/copy"
)

// CreateProject will create projects.
// It will create directories and files, copy templates and run scripts.
func CreateProject(label string, projects []string) error {
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
	if err = newSetup.init(); err != nil {
		return err
	}
	defer newSetup.stop()

	// Check if label is supported
	id, err := newSetup.isLabelSupported()
	if err != nil {
		return err
	}

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
	DatabaseName string
	Label        string
	ConfigDir    string
	InstallDir   string
	templatesDir string
	scriptsDir   string
	dbDir        string
	db           *sql.DB
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
func (setup *Setup) isLabelSupported() (int, error) {
	stmt, err := setup.db.Prepare("SELECT class_id FROM class_label WHERE label = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	var id int
	if err = stmt.QueryRow(setup.Label).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

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
func (project *Project) create(projectID int) error {
	// Set id and installDir
	project.ID = projectID
	project.InstallDir = project.Data.Owd + project.Name

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
	// Prepare statement for class folders
	stmtClass, err := project.Data.db.Prepare("SELECT target FROM class_folder WHERE class_id = ? AND template IS NULL")
	if err != nil {
		return err
	}
	defer stmtClass.Close()

	subFoldersClass, err := stmtClass.Query(project.ID)
	if err != nil {
		return err
	}
	defer subFoldersClass.Close()

	// Prepare statement for global folders
	stmtGlobal, err := project.Data.db.Prepare("SELECT target FROM global_folder WHERE template IS NULL")
	if err != nil {
		return err
	}
	defer stmtClass.Close()

	subFoldersGlobal, err := stmtGlobal.Query()
	if err != nil {
		return err
	}
	defer subFoldersGlobal.Close()

	// Create subfolders
	allSubFolders := []*sql.Rows{subFoldersClass, subFoldersGlobal}
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for _, subFolders := range allSubFolders {
		for subFolders.Next() {
			var subFolder string
			err = subFolders.Scan(&subFolder)
			if err != nil {
				return err
			}

			// Replace env variables
			subFolder = re.ReplaceAllString(subFolder, project.Name)

			// Create folder
			err = os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// createFiles queries all files from the database related to the projectId.
// Tries to create all of the files in the projectfolder.
// Returns error on failure. Returns nil on success.
func (project *Project) createFiles() error {
	// Prepare statement for class files
	stmtClass, err := project.Data.db.Prepare("SELECT target FROM class_file WHERE class_id = ? AND template IS NULL")
	if err != nil {
		return err
	}
	defer stmtClass.Close()

	filesClass, err := stmtClass.Query(project.ID)
	if err != nil {
		return err
	}
	defer filesClass.Close()

	// Prepare statement for global files
	stmtGlobal, err := project.Data.db.Prepare("SELECT target FROM global_file WHERE template IS NULL")
	if err != nil {
		return err
	}
	defer stmtClass.Close()

	filesGlobal, err := stmtGlobal.Query()
	if err != nil {
		return err
	}
	defer filesGlobal.Close()

	allFiles := []*sql.Rows{filesClass, filesGlobal}
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for _, files := range allFiles {
		for files.Next() {
			var file string
			err = files.Scan(&file)
			if err != nil {
				return err
			}

			// Replace env variables
			file = re.ReplaceAllString(file, project.Name)

			// Create file
			_, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// copyTemplates queries all templates from the database related to the projectId.
// Tries to copy all of the templates into the projectfolder.
// Returns error on failure. Returns nil on success.
func (project *Project) copyTemplates() error {
	// Prepare statement for class folders
	stmt, err := project.Data.db.Prepare("SELECT target, template FROM class_folder WHERE class_id = ? AND template IS NOT NULL")
	if err != nil {
		return err
	}
	defer stmt.Close()

	subFoldersClass, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer subFoldersClass.Close()

	// Prepare statement for global folders
	if stmt, err = project.Data.db.Prepare("SELECT target, template FROM global_folder WHERE template IS NOT NULL"); err != nil {
		return err
	}
	subFoldersGlobal, err := stmt.Query()
	if err != nil {
		return err
	}
	defer subFoldersGlobal.Close()

	// Prepare statement for class files
	if stmt, err = project.Data.db.Prepare("SELECT target, template FROM class_file WHERE class_id = ? AND template IS NOT NULL"); err != nil {
		return err
	}
	filesClass, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer filesClass.Close()

	// Prepare statement for global files
	if stmt, err = project.Data.db.Prepare("SELECT target, template FROM global_file WHERE template IS NOT NULL"); err != nil {
		return err
	}
	filesGlobal, err := stmt.Query()
	if err != nil {
		return err
	}
	defer filesGlobal.Close()

	templatesData := []*sql.Rows{subFoldersClass, subFoldersGlobal, filesClass, filesGlobal}

	for _, templateData := range templatesData {
		for templateData.Next() {
			var target, template string
			err = templateData.Scan(&target, &template)
			if err != nil {
				return err
			}

			template = project.Data.templatesDir + template
			err := copy.Copy(template, target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// runScripts queries all scripts from the database related to the projectId.
// Tries to execute all scripts.
// Returns error on failure. Returns nil on success.
func (project *Project) runScripts() error {
	// Prepare statement for class scripts
	stmt, err := project.Data.db.Prepare("SELECT name, run_as_sudo FROM class_script WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	classScripts, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer classScripts.Close()

	// Prepare statement for global scripts
	if stmt, err = project.Data.db.Prepare("SELECT name, run_as_sudo FROM global_script"); err != nil {
		return err
	}

	globalScripts, err := stmt.Query()
	if err != nil {
		return err
	}
	defer globalScripts.Close()

	allScripts := []*sql.Rows{classScripts, globalScripts}

	// Create scripts
	for _, scripts := range allScripts {
		for scripts.Next() {
			var script string
			var runAsSudo bool
			err = scripts.Scan(&script, &runAsSudo)
			if err != nil {
				return err
			}

			script = project.Data.scriptsDir + script

			if runAsSudo {
				script = "sudo " + script
			}

			cmd := exec.Command(script)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
