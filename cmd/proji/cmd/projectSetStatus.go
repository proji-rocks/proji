package cmd

import (
	"fmt"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var projectSetStatusCmd = &cobra.Command{
	Use:   "status STATUS PROJECT-ID",
	Short: "Set a new status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing status or project-ID")
		}

		status := strings.ToLower(args[0])
		projectID, err := helper.StrToUInt(args[1])
		if err != nil {
			return err
		}

		if err := setStatus(projectID, status); err != nil {
			fmt.Printf("Setting status '%s' for project %d failed: %v\n", status, projectID, err)
			return err
		}
		fmt.Printf("Status '%s' was successfully set for project %d.\n", status, projectID)
		return nil
	},
}

func init() {
	projectSetCmd.AddCommand(projectSetStatusCmd)
}

func setStatus(projectID uint, statusTitle string) error {
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	// Load and validate status
	status, err := s.LoadStatusByTitle(statusTitle)
	if err != nil {
		return err
	}
	return s.UpdateProjectStatus(projectID, status.ID)
}
