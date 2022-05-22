package cmd

import (
	"fmt"

	"github.com/brokeyourbike/mountefi/disk"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumesCmd)
}

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "List volumes",
	RunE: func(cmd *cobra.Command, args []string) error {
		volumes, err := disk.GetMountedVolumes()
		if err != nil {
			return err
		}

		for _, v := range volumes {
			info, err := disk.GetVolumeInfo(v)
			if err != nil {
				fmt.Printf("err %v", err)
				continue
			}
			fmt.Printf("%s (%s)\n", v, info.DeviceIdentifier)
		}

		return nil
	},
}
