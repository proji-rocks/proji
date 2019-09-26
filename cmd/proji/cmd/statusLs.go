package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var statusLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listStatuses()
	},
}

func init() {
	statusCmd.AddCommand(statusLsCmd)
}

func listStatuses() error {
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

	statuses, err := s.ListAvailableStatuses()
	if err != nil {
		return err
	}

	// Table header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Comment"})

	// Fill table
	for _, status := range statuses {
		t.AppendRow([]interface{}{status.ID, status.Title, status.Comment})
	}

	// Print the table
	t.Render()
	return nil
}
