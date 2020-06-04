package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/cobra"
)

var removeAllStatuses, forceRemoveStatuses bool

var statusRmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more statuses",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect statuses that will be removed
		var statuses []*item.Status

		if removeAllStatuses {
			var err error
			statuses, err = projiEnv.Svc.LoadAllStatuses()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing status-id")
			}

			for _, idStr := range args {
				id, err := helper.StrToUInt(idStr)
				if err != nil {
					return err
				}
				status, err := projiEnv.Svc.LoadStatus(id)
				if err != nil {
					return err
				}
				statuses = append(statuses, status)
			}
		}

		// Remove the statuses
		for _, status := range statuses {
			// Skip default statuses
			if status.IsDefault {
				continue
			}
			// Ask for confirmation if force flag was not passed
			if !forceRemoveClasses {
				if !helper.WantTo(
					fmt.Sprintf("Do you really want to remove status '%s (%d)'?", status.Title, status.ID),
				) {
					continue
				}
			}
			err := projiEnv.Svc.RemoveStatus(status.ID)
			if err != nil {
				fmt.Printf("> Removing status '%d' failed: %v\n", status.ID, err)
				return err
			}
			fmt.Printf("> Status '%d' was successfully removed\n", status.ID)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusRmCmd)
	statusRmCmd.Flags().BoolVarP(&removeAllStatuses, "all", "a", false, "Remove all statuses")
	classRmCmd.Flags().BoolVarP(&forceRemoveStatuses, "force", "f", false, "Don't ask for confirmation")
}
