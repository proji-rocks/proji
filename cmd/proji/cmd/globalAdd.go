package cmd

import (
	"fmt"
	"strings"

	"github.com/nikoksr/proji/internal/app/proji/global"
	"github.com/spf13/cobra"
)

var globalAddType string

var globalAddCmd = &cobra.Command{
	Use:   "add TARGET [TEMPLATE]",
	Short: "add a new global",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing global name")
		}
		if len(args) > 2 {
			return fmt.Errorf("too many arguments were given")
		}
		if err := isTypeValid(globalAddType); err != nil {
			return err
		}
		return global.AddGlobal(globalAddType, args)
	},
}

func init() {
	globalCmd.AddCommand(globalAddCmd)
	globalAddCmd.PersistentFlags().StringVarP(&globalAddType, "type", "t", "", "type of global - folder, file or script")
	globalAddCmd.MarkPersistentFlagRequired("type")
}

func isTypeValid(typeName string) error {
	switch strings.ToLower(typeName) {
	case "folder", "file", "script":
		return nil
	default:
		return fmt.Errorf("no valid type given")
	}
}
