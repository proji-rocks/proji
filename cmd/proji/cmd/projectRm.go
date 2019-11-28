package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var rmAll bool

var rmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if rmAll {
			err := removeAllProjects(projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing of all projects failed: %v\n", err)
				return err
			}
			fmt.Println("> All projects were successfully removed")
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("missing project id")
		}

		for _, idStr := range args {
			id, err := helper.StrToUInt(idStr)
			if err != nil {
				return err
			}

			err = removeProject(id, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing project '%d' failed: %v\n", id, err)
				continue
			}
			fmt.Printf("> Project '%d' was successfully removed\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&rmAll, "all", "a", false, "Remove all projects")
}

func removeProject(projectID uint, svc storage.Service) error {
	// Check if project exists
	_, err := svc.LoadProject(projectID)
	if err != nil {
		return err
	}
	return svc.RemoveProject(projectID)
}

func removeAllProjects(svc storage.Service) error {
	projects, err := svc.LoadAllProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		err = svc.RemoveProject(project.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
