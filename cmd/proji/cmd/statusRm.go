package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/cobra"
)

var rmAllStatuses bool

var statusRmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more statuses",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect statuses that will be removed
		var statuses []*item.Status

		if rmAllStatuses {
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
			if status.IsDefault {
				continue
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
	statusRmCmd.Flags().BoolVarP(&rmAllStatuses, "all", "a", false, "Remove all statuses")
}
