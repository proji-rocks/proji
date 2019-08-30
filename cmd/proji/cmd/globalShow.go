package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/global"
	"github.com/spf13/cobra"
)

var globalShowType string

var globalShowCmd = &cobra.Command{
	Use:   "show GLOBAL-ID",
	Short: "show detailed information about a global",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing global id")
		}
		if len(args) > 2 {
			return fmt.Errorf("too many global ids given")
		}
		if err := isTypeValid(globalShowType); err != nil {
			return err
		}

		return global.Show(globalShowType, args[0])
	},
}

func init() {
	globalCmd.AddCommand(globalShowCmd)
	globalShowCmd.PersistentFlags().StringVarP(&globalShowType, "type", "t", "", "type of global - folder, file or script")
	globalShowCmd.MarkPersistentFlagRequired("type")
}
