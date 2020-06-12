package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var statusAddCmd = &cobra.Command{
	Use:        "add STATUS [STATUS...]",
	Short:      "Add one or more statuses",
	Deprecated: "support for project statuses will be dropped with v0.21.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing status")
		}

		for _, status := range args {
			status = strings.ToLower(status)
			comment, err := addStatus(status)
			if err != nil {
				fmt.Printf("> Adding status %s failed: %v\n", status, err)
				if err.Error() == "Status already exists" {
					if !helper.WantTo("> Do you want to update its comment?") {
						continue
					}
					err := replaceStatus(status, comment)
					if err != nil {
						fmt.Printf("> Updating comment %s failed: %v\n", status, err)
						continue
					}
					fmt.Printf("> Comment for status '%s' was successfully updated\n", status)
				}
				continue
			}
			fmt.Printf("> Status '%s' was successfully created\n", status)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusAddCmd)
}

func addStatus(title string) (string, error) {
	// Create status and set status
	var status item.Status
	status.Title = title

	// Get a comment describing the status
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> Comment: ")
	comment, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	status.Comment = strings.Trim(comment, "\n")
	return status.Comment, projiEnv.Svc.SaveStatus(&status)
}

func replaceStatus(title, comment string) error {
	id, err := projiEnv.Svc.LoadStatusID(title)
	if err != nil {
		return err
	}
	return projiEnv.Svc.UpdateStatus(id, title, comment)
}
