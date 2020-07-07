//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/storage"

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
			return err
		}

		// Load class once for all projects
		class, err := session.StorageService.LoadClass(label)
		if err != nil {
			return err
		}

		for _, projectName := range projectNames {
			fmt.Printf("\n> Creating project %s\n", projectName)

			// Try to create the project
			projectPath := filepath.Join(workingDirectory, projectName)
			err := createProject(projectName, projectPath, class)
			if err == nil {
				fmt.Printf("> Project %s was successfully created\n", projectName)
				continue
			}

			// Print error message
			fmt.Printf(" > Failed: %v\n", err)

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
			err = replaceProject(projectName, projectPath, class)
			if err != nil {
				fmt.Printf("> Replacing project %s failed: %v\n", projectName, err)
				continue
			}
			fmt.Printf("> Project %s was successfully replaced\n", projectName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// createProject is a small wrapper function which takes a project name, path and its associated class,
// creates the project directory and tries to save it to storage.
func createProject(name, path string, class *models.Class) error {
	project := models.NewProject(name, path, class)
	err := project.Create(session.Config.BasePath)
	if err != nil {
		return err
	}
	err = session.StorageService.SaveProject(project)
	if err != nil {
		return err
	}
	return nil
}

// replaceProject should usually be executed after a attempt to create a new project failed with an ProjectExistsError.
// It will remove the given project from storage and save the new one, effectively replacing everything that's
// associated with the given project path.
func replaceProject(name, path string, class *models.Class) error {
	err := session.StorageService.RemoveProject(path)
	if err != nil {
		return err
	}
	project := models.NewProject(name, path, class)
	return session.StorageService.SaveProject(project)
}
