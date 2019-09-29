package cmd

import (
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
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

// cleanProjects cleans all projects.
func cleanProjects() error {
	// Setup storage service
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	projects, err := s.LoadAllProjects()
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
				if err := s.RemoveProject(project.ID); err != nil {
					return err
				}
				continue
			}
			if !statusGood {
				// Create status 'unknown'
				status := storage.Status{
					Title:   "unknown",
					Comment: "The state of this project is unknown.",
				}
				// Try adding the status to the storage
				if err := s.SaveStatus(&status); err != nil {
					if err.Error() != "Status '"+status.Title+"' already exists" {
						return err
					}
				}
				// Close and reconnect to load new data
				s.Close()
				s = nil
				s, err = sqlite.New(sqlitePath)
				if err != nil {
					return err
				}

				// Get its ID
				status.ID, err = s.LoadStatusID(status.Title)
				if err != nil {
					return err
				}
				// Update projects status to unknown
				if err = s.UpdateProjectStatus(project.ID, status.ID); err != nil {
					return err
				}
			}
			continue
		}
	}
	return nil
}
