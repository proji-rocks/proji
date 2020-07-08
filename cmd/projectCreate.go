//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/messages"

	"github.com/nikoksr/proji/storage"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"
	"github.com/nikoksr/proji/util"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create LABEL NAME [NAME...]",
	Short: "Create one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("at least one label and name have to be given")
		}
		label := args[0]
		projectNames := args[1:]

		// Get current working directory
		workingDirectory, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get working directory")
		}

		// Load package once for all projects
		pkg, err := session.StorageService.LoadPackage(label)
		if err != nil {
			return errors.Wrap(err, "failed to load package")
		}

		for _, projectName := range projectNames {
			messages.Info("creating project %s", projectName)

			// Try to create the project
			projectPath := filepath.Join(workingDirectory, projectName)
			err := createProject(projectName, projectPath, pkg)
			if err == nil {
				messages.Success("successfully created project %s", projectName)
				continue
			}

			// Print error message
			messages.Warning("failed to create project, %s", projectName, err.Error())

			// Check if error is because of a project is already associated with this path. Continue loop if so.
			_, projectExists := err.(*storage.ProjectExistsError)
			if !projectExists {
				continue
			}

			// Continue if use doesn't want to replace the project.
			if !util.WantTo("> Do you want to replace it?") {
				continue
			}

			// Try to replace the project
			err = replaceProject(projectName, projectPath, pkg)
			if err != nil {
				messages.Warning("failed to replace project %s, %s", projectName, err.Error())
			} else {
				messages.Success("successfully replaced project %s", projectName)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// createProject is a small wrapper function which takes a project name, path and its associated package,
// creates the project directory and tries to save it to storage.
func createProject(name, path string, pkg *models.Package) error {
	project := models.NewProject(name, path, pkg)
	err := project.Create(session.Config.BasePath)
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}
	err = session.StorageService.SaveProject(project)
	if err != nil {
		return errors.Wrap(err, "failed to save project")
	}
	return nil
}

// replaceProject should usually be executed after a attempt to create a new project failed with an ProjectExistsError.
// It will remove the given project from storage and save the new one, effectively replacing everything that's
// associated with the given project path.
func replaceProject(name, path string, pkg *models.Package) error {
	err := session.StorageService.RemoveProject(path)
	if err != nil {
		return errors.Wrap(err, "failed to remove project")
	}
	project := models.NewProject(name, path, pkg)
	err = session.StorageService.SaveProject(project)
	if err != nil {
		return errors.Wrap(err, "failed to save project")
	}
	return nil
}
