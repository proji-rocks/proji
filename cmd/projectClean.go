package cmd

import (
	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/internal/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type projectCleanCommand struct {
	cmd *cobra.Command
}

func newProjectCleanCommand() *projectCleanCommand {
	var cmd = &cobra.Command{
		Use:                   "clean",
		Short:                 "Clean up projects",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cleanProjects()
		},
	}
	return &projectCleanCommand{cmd: cmd}
}

func cleanProjects() error {
	projects, err := session.projectService.LoadProjectList()
	if err != nil {
		return errors.Wrap(err, "failed to load all projects")
	}

	for _, project := range projects {
		// Check path
		if util.DoesPathExist(project.Path) {
			continue
		}
		// Remove the project
		err := session.projectService.RemoveProject(project.Path)
		if err != nil {
			message.Warningf("failed to remove project with path %s, %v", project.Path, err)
		}
	}
	return nil
}
