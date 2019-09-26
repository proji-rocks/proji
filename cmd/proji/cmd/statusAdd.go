package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var statusAddCmd = &cobra.Command{
	Use:   "add STATUS [STATUS...]",
	Short: "Add one or more statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing status")
		}

		for _, status := range args {
			status = strings.ToLower(status)
			comment, err := addStatus(status)
			if err != nil {
				fmt.Printf("Adding status %s failed: %v\n", status, err)
				if err.Error() == "Status already exists" {
					if !helper.WantTo("Do you want to update its comment?") {
						continue
					}
					if err := replaceStatus(status, comment); err != nil {
						fmt.Printf("Updating comment %s failed: %v\n", status, err)
						continue
					}
					fmt.Printf("Comment for status '%s' was successfully updated.\n", status)
				}
				continue
			}
			fmt.Printf("Status '%s' was successfully created.\n", status)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusAddCmd)
}

func addStatus(statusTitle string) (string, error) {
	// Setup storage
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return "", err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return "", err
	}
	defer s.Close()

	// Create status and set status
	var status storage.Status
	status.Title = statusTitle

	// Get a comment describing the status
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Comment: ")
	comment, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	status.Comment = strings.Trim(comment, "\n")
	return status.Comment, s.AddStatus(&status)
}

func replaceStatus(status, comment string) error {
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

	id, err := s.LoadStatusID(status)
	if err != nil {
		return err
	}
	return s.UpdateStatus(&storage.Status{ID: id, Title: status, Comment: comment})
}
