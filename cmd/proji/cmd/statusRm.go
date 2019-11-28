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
			return fmt.Errorf("missing status-ID")
		}

		for _, status := range args {
			err := removeStatus(status, projiEnv.Svc)
			if err != nil {
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

func removeStatus(id string, svc storage.Service) error {
	statusID, err := helper.StrToUInt(id)
	if err != nil {
		return err
	}

	status, err := svc.LoadStatus(statusID)
	if err != nil {
		return err
	}

	if status.IsDefault {
		return fmt.Errorf("default statuses can not be removed")
	}
	return svc.RemoveStatus(statusID)
}

func removeAllStatuses(svc storage.Service) error {
	statuses, err := svc.LoadAllStatuses()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		if status.IsDefault {
			continue
		}
		err = svc.RemoveStatus(status.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
