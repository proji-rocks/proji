package models

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
	"gorm.io/gorm"
)

// Project represents a project that was created by proji. It holds tags for gorm and toml defining its storage and
// export/import behaviour.
type Project struct {
	ID        uint           `gorm:"primarykey" toml:"-"`
	CreatedAt time.Time      `toml:"-"`
	UpdatedAt time.Time      `gorm:"index:idx_unq_project_path_deletedat,unique;" toml:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" toml:"-"`
	Path      string         `gorm:"index:idx_unq_project_path_deletedat,unique;not null" toml:"label"`
	Class     *Class         `gorm:"ForeignKey:ID;References:ID"`
}

// NewProject returns a new project.
func NewProject(name, path string, class *Class) *Project {
	return &Project{
		Path:  path,
		Class: class,
	}
}

// Create starts the creation of a project.
func (p *Project) Create(cwd, configPath string) (err error) {
	err = p.createProjectFolder()
	if err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	err = os.Chdir(p.Path)
	if err != nil {
		return err
	}

	// Append a slash if not exists. Out of convenience.
	if cwd[:len(cwd)-1] != "/" {
		cwd += "/"
	}
	defer func() {
		newErr := os.Chdir(cwd)
		if newErr != nil {
			err = newErr
		}
	}()

	err = p.preRunPlugins(configPath)
	if err != nil {
		return err
	}

	err = p.createFilesAndFolders(configPath)
	if err != nil {
		return err
	}

	return p.postRunPlugins(configPath)
}

// createProjectFolder tries to create the main project folder.
func (p *Project) createProjectFolder() error {
	return os.Mkdir(p.Path, os.ModePerm)
}

func (p *Project) createFilesAndFolders(configPath string) error {
	templatePath := filepath.Join(configPath, "/templates/")
	for _, template := range p.Class.Templates {
		if len(template.Path) > 0 {
			// Copy template file or folder
			err := copy.Copy(filepath.Join(templatePath, "/", template.Path), template.Destination)
			if err != nil {
				return err
			}
		}
		if template.IsFile {
			// Create file
			_, err := os.Create(template.Destination)
			if err != nil {
				return err
			}
		} else {
			// Create folder
			err := os.MkdirAll(template.Destination, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Project) preRunPlugins(configPath string) error {
	for _, plugin := range p.Class.Plugins {
		if plugin.ExecNumber >= 0 {
			continue
		}
		pluginPath := filepath.Join(configPath, "/plugins/", plugin.Path)
		err := runPlugin(pluginPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) postRunPlugins(configPath string) error {
	for _, plugin := range p.Class.Plugins {
		if plugin.ExecNumber <= 0 {
			continue
		}
		pluginPath := filepath.Join(configPath, "/plugins/", plugin.Path)
		err := runPlugin(pluginPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func runPlugin(pluginPath string) error {
	cmd := exec.Command(pluginPath)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
