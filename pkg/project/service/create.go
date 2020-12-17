package projectservice

import (
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/domain"
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
		e := os.Chdir(workingDirectory)
		if e != nil {
			if err != nil {
				err = errors.Wrap(err, e.Error())
			} else {
				err = e
			}
		}
	}()

	// Run plugins before creation of sub-folders and files
	pluginsRootPath := filepath.Join(configRootPath, "plugins")
	err = preRunPlugins(pluginsRootPath, project.Package.Plugins)
	if err != nil {
		return err
	}

	// Create sub-folders and files
	templatesRootPath := filepath.Join(configRootPath, "templates")
	err = ps.templateEngine.CreateFilesInProjectFolder(templatesRootPath, project.Path, project.Package.Templates)
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
