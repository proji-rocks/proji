package cmd

import (
	"fmt"
	"strconv"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var projectRmCmd = &cobra.Command{
	Use:   "rm PROJECT-ID [PROJECT-ID...]",
	Short: "remove existing projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing project id")
		}

		for _, idStr := range args {
			// Parse the input
			id64, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return err
			}
			id := uint(id64)

			if err := RemoveProject(id); err != nil {
				fmt.Printf("Removing project with id %d failed: %v\n", id, err)
				continue
			}
			fmt.Printf("Project with id %d was successfully removed.\n", id)
		}
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectRmCmd)
}

// RemoveProject removes a project from storage.
func RemoveProject(projectID uint) error {
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
	return s.UntrackProject(projectID)
}
