package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var classRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}

		for _, name := range args {
			if err := removeClass(name, projiEnv.Svc); err != nil {
				fmt.Printf("> Removing '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("> '%s' was successfully removed\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
}

func removeClass(label string, svc storage.Service) error {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}
	return svc.RemoveClass(classID)
}
