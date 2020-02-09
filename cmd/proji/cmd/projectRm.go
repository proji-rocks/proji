package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/cobra"
)

var rmAll bool

var rmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect projects that will be removed
		var projects []*item.Project

		if rmAll {
			var err error
			projects, err = projiEnv.Svc.LoadAllProjects()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing project id")
			}

			for _, idStr := range args {
				id, err := helper.StrToUInt(idStr)
				if err != nil {
					return err
				}
				project, err := projiEnv.Svc.LoadProject(id)
				if err != nil {
					return err
				}
				projects = append(projects, project)
			}
		}

		// Remove the projects
		for _, project := range projects {
			err := projiEnv.Svc.RemoveProject(project.ID)
			if err != nil {
				fmt.Printf("> Removing project '%d' failed: %v\n", project.ID, err)
				return err
			}
			fmt.Printf("> Project '%d' was successfully removed\n", project.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&rmAll, "all", "a", false, "Remove all projects")
}
