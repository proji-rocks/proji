package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
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
			if err := removeClass(name); err != nil {
				fmt.Printf("Removing '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("'%s' was successfully removed.\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
}

func removeClass(label string) error {
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
	classID, err := s.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}
	return s.RemoveClass(classID)
}
