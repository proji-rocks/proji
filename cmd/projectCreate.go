package cmd

import (
	"os"
	"path/filepath"

	projectstore "github.com/nikoksr/proji/pkg/project/store"

	"github.com/nikoksr/proji/pkg/domain"

	"github.com/nikoksr/proji/internal/message"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/internal/util"
	"github.com/spf13/cobra"
)

type projectCreateCommand struct {
	cmd *cobra.Command
}

func newProjectCreateCommand() *projectCreateCommand {
	cmd := &cobra.Command{
		Use:                   "create LABEL NAME [NAME...]",
		Short:                 "Create one or more projects",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			label := args[0]
			projectNames := args[1:]

			// Get current working directory
			workingDirectory, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to get working directory")
			}

			// Load package once for all projects
			pkg, err := session.packageService.LoadPackage(true, label)
			if err != nil {
				return errors.Wrap(err, "failed to load package")
			}

			for _, projectName := range projectNames {
				message.Infof("creating project %s", projectName)

				// Try to create the project
				projectPath := filepath.Join(workingDirectory, projectName)
				err := createProject(projectName, projectPath, pkg)
				if err == nil {
					message.Successf("successfully created project %s", projectName)
					continue
				}

				// Print error message
				message.Warningf("failed to create project, %s", projectName, err.Error())

				// Check if error is because of a project is already associated with this path. Continue loop if so.
				if errors.Is(err, projectstore.ErrProjectExists) {
					continue
				}

				// Continue if use doesn't want to replace the project.
				if !util.WantTo("> Do you want to replace it?") {
					continue
				}

				// Try to replace the project
				err = replaceProject(projectName, projectPath, pkg)
				if err != nil {
					message.Warningf("failed to replace project %s, %s", projectName, err.Error())
				} else {
					message.Successf("successfully replaced project %s", projectName)
				}
			}
			return nil
		},
	}

	return &projectCreateCommand{cmd: cmd}
}

// createProject is a small wrapper function which takes a project name, path and its associated package,
// creates the project directory and tries to save it to storage.
func createProject(name, path string, pkg *domain.Package) error {
	project := domain.NewProject(name, path, pkg)
	err := session.projectService.CreateProject(session.config.BasePath, project)
	if err != nil {
		return errors.Wrap(err, "create project")
	}
	err = session.projectService.StoreProject(project)
	if err != nil {
		return errors.Wrap(err, "save project")
	}
	return nil
}

// replaceProject should usually be executed after a attempt to create a new project failed with an ErrProjectExists.
// It will remove the given project from storage and save the new one, effectively replacing everything that's
// associated with the given project path.
func replaceProject(name, path string, pkg *domain.Package) error {
	err := session.projectService.RemoveProject(path)
	if err != nil {
		return errors.Wrap(err, "remove project")
	}
	project := domain.NewProject(name, path, pkg)
	err = session.projectService.StoreProject(project)
	if err != nil {
		return errors.Wrap(err, "save project")
	}
	return nil
}
