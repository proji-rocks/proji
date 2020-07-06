//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/storage/models"
	"github.com/nikoksr/proji/util"
	"github.com/spf13/cobra"
)

var removeAllProjects, forceRemoveProjects bool

var rmCmd = &cobra.Command{
	Use:   "rm PATH [PATH...]",
	Short: "Remove one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect projects that will be removed
		var projects []*models.Project

		if removeAllProjects {
			var err error
			projects, err = projiEnv.StorageService.LoadAllProjects()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing project paths")
			}

			for _, path := range args {
				project, err := projiEnv.StorageService.LoadProject(path)
				if err != nil {
					return err
				}
				projects = append(projects, project)
			}
		}

		// Remove the projects
		for _, project := range projects {
			// Ask for confirmation if force flag was not passed
			if !forceRemoveProjects {
				if !util.WantTo(
					fmt.Sprintf("Do you really want to remove project '%s (%d)'?", project.Name, project.ID),
				) {
					continue
				}
			}
			err := projiEnv.StorageService.RemoveProject(project.Path)
			if err != nil {
				fmt.Printf("> Removing project '%s' failed: %v\n", project.Path, err)
				return err
			}
			fmt.Printf("> Project '%s' was successfully removed\n", project.Path)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&removeAllClasses, "all", "a", false, "Remove all projects")
	rmCmd.Flags().BoolVarP(&forceRemoveProjects, "force", "f", false, "Don't ask for confirmation")
}
