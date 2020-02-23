package item

import (
	"os"
	"os/exec"
	"path/filepath"
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
	err := proj.createProjectFolder()
	if err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	err = os.Chdir(proj.Name)
	if err != nil {
		return err
	}

	// Append a slash if not exists. Out of convenience.
	if cwd[:len(cwd)-1] != "/" {
		cwd += "/"
	}
	defer os.Chdir(cwd)

	err = proj.preRunScripts(configPath)
	if err != nil {
		return err
	}

	err = proj.createSubFolders()
	if err != nil {
		return err
	}

	err = proj.createFiles()
	if err != nil {
		return err
	}

	err = proj.copyTemplates(configPath)
	if err != nil {
		return err
	}
	return proj.postRunScripts(configPath)
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
		err := os.MkdirAll(folder.Destination, os.ModePerm)
		if err != nil {
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
		_, err := os.OpenFile(file.Destination, os.O_RDONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (proj *Project) copyTemplates(configPath string) error {
	re := regexp.MustCompile(`__PROJECT_NAME__`)
	templatePath := filepath.Join(configPath, "/templates/")

	for _, folder := range proj.Class.Folders {
		if len(folder.Template) < 1 {
			continue
		}

		// Replace keyword with project name
		folder.Destination = re.ReplaceAllString(folder.Destination, proj.Name)
		err := copy.Copy(filepath.Join(templatePath, "/", folder.Template), folder.Destination)
		if err != nil {
			return err
		}
	}

	for _, file := range proj.Class.Files {
		if len(file.Template) < 1 {
			continue
		}

		// Replace keyword with project name
		file.Destination = re.ReplaceAllString(file.Destination, proj.Name)
		err := copy.Copy(filepath.Join(templatePath, "/", file.Template), file.Destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (proj *Project) preRunScripts(configPath string) error {
	return proj.runScripts("pre", configPath)
}

func (proj *Project) postRunScripts(configPath string) error {
	return proj.runScripts("post", configPath)
}

func (proj *Project) runScripts(scriptType, configPath string) error {
	for _, script := range proj.Class.Scripts {
		if script.Type != scriptType {
			continue
		}

		scriptPath := filepath.Join(configPath, "/scripts/", script.Name)

		if script.RunAsSudo {
			scriptPath = "sudo " + scriptPath
		}

		re := regexp.MustCompile(`__PROJECT_NAME__`)
		for idx, arg := range script.Args {
			script.Args[idx] = re.ReplaceAllString(arg, proj.Name)
		}

		cmd := exec.Command(scriptPath, script.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
