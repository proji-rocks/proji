package cmd

import (
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleanProjects()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanProjects() error {
	projects, err := projiEnv.Svc.LoadAllProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		// Check path
		pathGood := true
		if !helper.DoesPathExist(project.InstallPath) {
			pathGood = false
		}

		// Check status
		statusGood := true
		if project.Status.Title == "" {
			statusGood = false
		}

		// Overall health
		overallGood := pathGood && statusGood

		if !overallGood {
			if !pathGood {
				// Remove the project
				err := projiEnv.Svc.RemoveProject(project.ID)
				if err != nil {
					return err
				}
				continue
			}
			if !statusGood {
				// Update projects status to unknown (ID 5 in storage)
				err = projiEnv.Svc.UpdateProjectStatus(project.ID, 5)
				if err != nil {
					return err
				}
				continue
			}
		}
	}
	return nil
}
