package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/messages"

	"github.com/nikoksr/proji/storage/models"
	"github.com/nikoksr/proji/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type projectRemoveCommand struct {
	cmd *cobra.Command
}

func newProjectRemoveCommand() *projectRemoveCommand {
	var removeAllProjects, forceRemoveProjects bool

	var cmd = &cobra.Command{
		Use:   "rm PATH [PATH...]",
		Short: "Remove one or more projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Collect projects that will be removed
			var projects []*models.Project

			if removeAllProjects {
				var err error
				projects, err = activeSession.storageService.LoadProjects()
				if err != nil {
					return errors.Wrap(err, "failed to load all projects")
				}
			} else {
				if len(args) < 1 {
					return fmt.Errorf("missing project paths")
				}

				for _, path := range args {
					project, err := activeSession.storageService.LoadProject(path)
					if err != nil {
						return errors.Wrap(err, "failed to load project")
					}
					projects = append(projects, project)
				}
			}

			// Remove the projects
			for _, project := range projects {
				// Ask for confirmation if force flag was not passed
				if !forceRemoveProjects {
					if !util.WantTo(
						fmt.Sprintf("Do you really want to remove the path %s from your projects?", project.Path),
					) {
						continue
					}
				}
				err := activeSession.storageService.RemoveProject(project.Path)
				if err != nil {
					messages.Warning("failed to remove project %s, %s", project.Path, err.Error())
					continue
				}
				messages.Success("successfully removed project %s", project.Path)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&removeAllProjects, "all", "a", false, "Remove all projects")
	cmd.Flags().BoolVarP(&forceRemoveProjects, "force", "f", false, "Don't ask for confirmation")
	return &projectRemoveCommand{cmd: cmd}
}
