package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var statusRmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing status-ID")
		}

		for _, status := range args {
			if err := removeStatus(status); err != nil {
				fmt.Printf("Removing status %s failed: %v\n", status, err)
			}
			fmt.Printf("Status '%s' was successfully removed.\n", status)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusRmCmd)
}

func removeStatus(status string) error {
	// Setup storage
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	statusID, err := helper.StrToUInt(status)
	if err != nil {
		return err
	}
	return s.RemoveStatus(statusID)
}
