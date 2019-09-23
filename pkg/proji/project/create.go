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
	if err = proj.new(id); err != nil {
		return fmt.Errorf("could not create project %s: %v", project, err)
	}

	return nil
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

			template = project.Data.ConfigDir + "/templates/" + template
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

		script = project.Data.ConfigDir + "/scripts/" + script

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
