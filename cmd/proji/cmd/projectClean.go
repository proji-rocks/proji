package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var dryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleanProjects(dryRun)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Don't auto clean. Only show dirty projects.")
}

// cleanProjects cleans all projects.
func cleanProjects(dryRun bool) error {
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

	projects, err := s.ListProjects()
	if err != nil {
		return err
	}

	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Project", "Path", "Status", "Overall"})

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
		if !dryRun && !overallGood {
			if !pathGood {
				// Remove the project
				if err := s.UntrackProject(project.ID); err != nil {
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
				if err := s.AddStatus(&status); err != nil {
					if err.Error() != "Status already exists" {
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

		// Dry-Run, print only
		path := boolToInfo(pathGood)
		status := boolToInfo(statusGood)
		overall := boolToInfo(overallGood)

		t.AppendRow([]interface{}{
			project.ID,
			path,
			status,
			overall,
		})
	}

	// Print the table
	if dryRun {
		t.Render()
	}
	return nil
}

func boolToInfo(good bool) string {
	if good {
		// return ""
		return "Good"
	}
	// return ""
	return "Bad"
}
