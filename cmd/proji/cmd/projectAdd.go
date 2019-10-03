package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add LABEL PATH STATUS",
	Short: "Add an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("Missing label, path or status")
		}

		path, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		if !helper.DoesPathExist(path) {
			return fmt.Errorf("path '%s' does not exist", path)
		}

		label := strings.ToLower(args[0])
		status := strings.ToLower(args[2])

		if err := addProject(label, path, status, projiEnv.Svc); err != nil {
			return err
		}
		fmt.Printf("Project '%s' was successfully added.\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addProject(label, path, statusTitle string, svc storage.Service) error {
	name := filepath.Base(path)
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}

	statusID, err := svc.LoadStatusID(statusTitle)
	if err != nil {
		return err
	}

	class, err := svc.LoadClass(classID)
	if err != nil {
		return err
	}

	var status *item.Status
	status, err = svc.LoadStatus(statusID)
	if err != nil {
		// Load status unknown
		status, err = svc.LoadStatus(5)
		if err != nil {
			return err
		}
	}

	proj, err := item.NewProject(0, name, path, class, status)
	if err != nil {
		return err
	}
	if err := svc.SaveProject(proj); err != nil {
		return err
	}
	return nil
}
