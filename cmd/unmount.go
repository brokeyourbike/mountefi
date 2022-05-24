package cmd

import (
	"fmt"

	"github.com/brokeyourbike/mountefi/disk"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unmountCmd)
}

var unmountCmd = &cobra.Command{
	Use:   "unmount <disk>",
	Short: "Unmount specified disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := disk.GetDiskInfo(args[0])
		if err != nil {
			return fmt.Errorf("can not get disk info for: %s", args[0])
		}

		if !info.IsEfi() {
			return fmt.Errorf("can not unmount %s. not an EFI", info.DeviceIdentifier)
		}

		if !info.IsMounted() {
			return fmt.Errorf("%s is not mounted", info.DeviceIdentifier)
		}

		out, err := disk.UnmountDisk(info)
		if err != nil {
			return fmt.Errorf("can not unmount disk %s", info.DeviceIdentifier)
		}

		fmt.Printf("disk unmounted: %s\n", out)
		return nil
	},
}
