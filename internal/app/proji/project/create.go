package project

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"

	// Import sqlite3 driver (see func (setup *Setup) Run() error)
	"github.com/mattn/go-sqlite3"
	"github.com/nikoksr/proji/internal/app/helper"

	"github.com/otiai10/copy"
)

// CreateProject will create a new project or return an error if the project already exists.
// It will create directories and files, copy templates and run scripts.
func CreateProject(label string, project string) error {
	configDir := helper.GetConfigDir()
	databaseName, ok := viper.Get("database.name").(string)

	if !ok {
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

	// Header
	fmt.Println(helper.ProjectHeader(project))
	proj := Project{Name: project, Data: &newSetup}
	// Track
	if err = proj.track(); err != nil {
		return fmt.Errorf("could not create project %s: %v", project, err)
	}
	// Create
	if err = proj.create(id); err != nil {
		return fmt.Errorf("could not create project %s: %v", project, err)
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

	if setup.Owd[:len(setup.Owd)-1] != "/" {
		setup.Owd += "/"
	}

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
	err = stmt.QueryRow(setup.Label).Scan(&id)
	return id, err
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

// createProjectFolder tries to create the main project folder.
// Returns an error on failure.
func (project *Project) createProjectFolder() error {
	return os.Mkdir(project.Name, os.ModePerm)
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

	subFolders, err := stmtClass.Query(project.ID)
	if err != nil {
		return err
	}
	defer subFolders.Close()
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for subFolders.Next() {
		var subFolder string
		if err = subFolders.Scan(&subFolder); err != nil {
			return err
		}

		// Replace variable with project name
		subFolder = re.ReplaceAllString(subFolder, project.Name)

		// Create folder
		if err = os.MkdirAll(subFolder, os.ModePerm); err != nil {
			return err
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

	files, err := stmtClass.Query(project.ID)
	if err != nil {
		return err
	}
	defer files.Close()
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for files.Next() {
		var file string
		if err = files.Scan(&file); err != nil {
			return err
		}

		// Replace variable with project name
		file = re.ReplaceAllString(file, project.Name)

		// Create file
		if _, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, os.ModePerm); err != nil {
			return err
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

	subFolders, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer subFolders.Close()

	// Prepare statement for class files
	if stmt, err = project.Data.db.Prepare("SELECT target, template FROM class_file WHERE class_id = ? AND template IS NOT NULL"); err != nil {
		return err
	}
	files, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer files.Close()

	templatesData := []*sql.Rows{subFolders, files}

	for _, templateData := range templatesData {
		for templateData.Next() {
			var target, template string
			if err = templateData.Scan(&target, &template); err != nil {
				return err
			}

			template = project.Data.templatesDir + template
			if err = copy.Copy(template, target); err != nil {
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

	scripts, err := stmt.Query(project.ID)
	if err != nil {
		return err
	}
	defer scripts.Close()

	// Create scripts
	for scripts.Next() {
		var script string
		var runAsSudo bool
		if err = scripts.Scan(&script, &runAsSudo); err != nil {
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
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
