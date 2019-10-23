package item

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/otiai10/copy"
)

// Project struct represents a proji project
type Project struct {
	ID          uint    // The project ID
	Name        string  // The project name
	InstallPath string  // The install path for the project
	Class       *Class  // The template class
	Status      *Status // The current project status
}

// NewProject returns a new project
func NewProject(projectID uint, name, installPath string, class *Class, status *Status) *Project {
	return &Project{
		ID:          projectID,
		Name:        name,
		InstallPath: installPath,
		Class:       class,
		Status:      status,
	}
}

// Create starts the creation of a project.
func (proj *Project) Create(cwd, configPath string) error {
	if err := proj.createProjectFolder(); err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	if err := os.Chdir(proj.Name); err != nil {
		return err
	}

	// Append a slash if not exists. Out of convenience.
	if cwd[:len(cwd)-1] != "/" {
		cwd += "/"
	}
	defer os.Chdir(cwd)

	if err := proj.createSubFolders(); err != nil {
		return err
	}
	if err := proj.createFiles(); err != nil {
		return err
	}
	if err := proj.copyTemplates(configPath); err != nil {
		return err
	}
	return proj.runScripts(configPath)
}

// createProjectFolder tries to create the main project folder.
func (proj *Project) createProjectFolder() error {
	return os.Mkdir(proj.Name, os.ModePerm)
}

func (proj *Project) createSubFolders() error {
	// Regex to replace keyword with project name
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for _, folder := range proj.Class.Folders {
		// Skip, folder has a template
		if len(folder.Template) > 0 {
			continue
		}

		// Replace keyword with project name
		folder.Destination = re.ReplaceAllString(folder.Destination, proj.Name)
		if err := os.MkdirAll(folder.Destination, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (proj *Project) createFiles() error {
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for _, file := range proj.Class.Files {
		// Skip, file has a template
		if len(file.Template) > 0 {
			continue
		}

		// Replace keyword with project name
		file.Destination = re.ReplaceAllString(file.Destination, proj.Name)
		if _, err := os.OpenFile(file.Destination, os.O_RDONLY|os.O_CREATE, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (proj *Project) copyTemplates(configPath string) error {
	re := regexp.MustCompile(`__PROJECT_NAME__`)

	for _, folder := range proj.Class.Folders {
		if len(folder.Template) < 1 {
			continue
		}

		// Replace keyword with project name
		folder.Destination = re.ReplaceAllString(folder.Destination, proj.Name)
		folder.Template = configPath + "templates/" + folder.Template
		if err := copy.Copy(folder.Template, folder.Destination); err != nil {
			return err
		}
	}

	for _, file := range proj.Class.Files {
		if len(file.Template) < 1 {
			continue
		}

		// Replace keyword with project name
		file.Destination = re.ReplaceAllString(file.Destination, proj.Name)
		file.Template = configPath + "templates/" + file.Template
		if err := copy.Copy(file.Template, file.Destination); err != nil {
			return err
		}
	}
	return nil
}

func (proj *Project) runScripts(configPath string) error {
	for _, script := range proj.Class.Scripts {
		scriptPath := configPath + "scripts/" + script.Name

		if script.RunAsSudo {
			scriptPath = "sudo " + scriptPath
		}

		cmd := exec.Command(scriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
