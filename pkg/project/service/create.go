package projectservice

import (
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
)

// Create starts the creation of a project.
func (ps projectService) CreateProject(configRootPath string, project *domain.Project) (err error) {
	// Create the root folder of the project.
	err = createProjectRootFolder(project.Path)
	if err != nil {
		return errors.Wrap(err, "create base folder")
	}

	// Get working directory. We will be changing directories, so we need to know, where we started from.
	workingDirectory, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get working directory")
	}

	// Change directory into the new project directory and defer chdir back to old cwd
	err = os.Chdir(project.Path)
	if err != nil {
		return err
	}

	defer func() {
		newErr := os.Chdir(workingDirectory)
		if newErr != nil {
			err = newErr
		}
	}()

	// Run plugins before creation of subfolders and files
	pluginsRootPath := filepath.Join(configRootPath, "plugins")
	err = preRunPlugins(pluginsRootPath, project.Package.Plugins)
	if err != nil {
		return err
	}

	// Create sub-folders and files
	err = createFilesAndFolders(configRootPath, project.Package.Templates)
	if err != nil {
		return err
	}

	// Run plugins after all folders and files have been created
	return postRunPlugins(pluginsRootPath, project.Package.Plugins)
}

// createProjectRootFolder tries to create the root project folder.
func createProjectRootFolder(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

func createFilesAndFolders(configRootPath string, templates []*domain.Template) error {
	baseTemplatesPath := filepath.Join(configRootPath, "/templates/")
	for _, template := range templates {
		if len(template.Path) > 0 {
			// Copy template file or folder
			err := copy.Copy(filepath.Join(baseTemplatesPath, template.Path), template.Destination)
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

func preRunPlugins(pluginsRootPath string, plugins []*domain.Plugin) error {
	for _, plugin := range plugins {
		if plugin.ExecNumber >= 0 {
			continue
		}
		// Plugin path is relative by default to make it shareable. We have to make it an absolute path here,
		// so that we can execute it.
		plugin.Path = filepath.Join(pluginsRootPath, plugin.Path)
		err := runPlugin(plugin.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func postRunPlugins(pluginsRootPath string, plugins []*domain.Plugin) error {
	for _, plugin := range plugins {
		if plugin.ExecNumber <= 0 {
			continue
		}
		// Plugin path is relative by default to make it shareable. We have to make it an absolute path here,
		// so that we can execute it.
		plugin.Path = filepath.Join(pluginsRootPath, plugin.Path)
		err := runPlugin(plugin.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func runPlugin(pluginPath string) error {
	luaState := lua.NewState()
	defer luaState.Close()
	return luaState.DoFile(pluginPath)
}
