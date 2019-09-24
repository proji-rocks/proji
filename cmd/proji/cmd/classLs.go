package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var classLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list existing classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ListClasses()
	},
}

func init() {
	classCmd.AddCommand(classLsCmd)
}

// ListClasses lists all classes available in the database
func ListClasses() error {
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

	classes, err := s.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		fmt.Println(class.Name)
	}

	return nil
}
