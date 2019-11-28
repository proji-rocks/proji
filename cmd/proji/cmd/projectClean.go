package cmd

import (
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleanProjects(projiEnv.Svc)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanProjects(svc storage.Service) error {
	projects, err := svc.LoadAllProjects()
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

		// If no dry run and overall healh is bad, than clean project
		if !overallGood {
			if !pathGood {
				// Remove the project
				err := svc.RemoveProject(project.ID)
				if err != nil {
					return err
				}
				continue
			}
			if !statusGood {
				// Update projects status to unknown (ID 5 in storage)
				err = svc.UpdateProjectStatus(project.ID, 5)
				if err != nil {
					return err
				}
				continue
			}
		}
	}
	return nil
}
