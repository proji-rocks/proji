package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classRmCmd = &cobra.Command{
	Use:   "rm NAME [NAME...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}

		for _, name := range args {
			if err := RemoveClass(name); err != nil {
				fmt.Printf("Removing class %s failed: %v\n", name, err)
				continue
			}
			fmt.Printf("Class %s was successfully removed.\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
}

// RemoveClass removes a class from storage.
func RemoveClass(name string) error {
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

	return s.RemoveClass(name)
}
