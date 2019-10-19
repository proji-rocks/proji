package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm ID [ID...]",
	Short: "Remove one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing project id")
		}

		for _, idStr := range args {
			id, err := helper.StrToUInt(idStr)
			if err != nil {
				return err
			}

			if err := removeProject(id, projiEnv.Svc); err != nil {
				fmt.Printf("> Removing project '%d' failed: %v\n", id, err)
				continue
			}
			fmt.Printf("> Project '%d' was successfully removed\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func removeProject(projectID uint, svc storage.Service) error {
	// Check if class exists
	if _, err := svc.LoadProject(projectID); err != nil {
		return err
	}
	return svc.RemoveProject(projectID)
}
