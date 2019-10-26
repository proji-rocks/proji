package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var rmAllStatuses bool

var statusRmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		if rmAllStatuses {
			err := removeAllStatuses(projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing of all statuses failed: %v\n", err)
				return err
			}
			fmt.Println("> All statuses were successfully removed")
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("Missing status-ID")
		}

		for _, status := range args {
			if err := removeStatus(status, projiEnv.Svc); err != nil {
				fmt.Printf("> Removing status %s failed: %v\n", status, err)
			}
			fmt.Printf("> Status '%s' was successfully removed\n", status)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusRmCmd)
	statusRmCmd.Flags().BoolVarP(&rmAllStatuses, "all", "a", false, "Remove all statuses")
}

func removeStatus(status string, svc storage.Service) error {
	statusID, err := helper.StrToUInt(status)
	if err != nil {
		return err
	}

	if statusID < 6 {
		return fmt.Errorf("statuses 1-5 can not be removed")
	}

	// Check if status exists
	if _, err := svc.LoadStatus(statusID); err != nil {
		return err
	}
	return svc.RemoveStatus(statusID)
}

func removeAllStatuses(svc storage.Service) error {
	statuses, err := svc.LoadAllStatuses()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		if status.ID < 6 {
			continue
		}
		err = svc.RemoveStatus(status.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
