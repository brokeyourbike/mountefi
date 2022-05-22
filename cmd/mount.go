package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mountCmd)
}

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
