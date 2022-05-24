package cmd

import (
	"fmt"

	"github.com/brokeyourbike/mountefi/disk"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

// sudoMountOsVersion is an instance of version.Version.
// It represents the version of MacOS that do not require sudo to mount.
var sudoMountOsVersion *version.Version

func init() {
	rootCmd.AddCommand(mountCmd)
	sudoMountOsVersion = version.Must(version.NewVersion("10.13.6"))
}

var mountCmd = &cobra.Command{
	Use:   "mount <disk>",
	Short: "Mount specified disk",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := disk.GetDiskInfo(args[0])
		if err != nil {
			return fmt.Errorf("can not get disk info for: %s", args[0])
		}

		list, err := disk.GetDiskList()
		if err != nil {
			return fmt.Errorf("can not get disks list. %v", err)
		}

		disks := disk.NewDisks(list.AllDisks)
		disks.Update()

		parent, err := disks.FindParentFor(info)
		if err != nil {
			return fmt.Errorf("can not get parent. %v", err)
		}

		fmt.Printf("Found parent: %v\n", parent.DeviceIdentifier)

		efi, err := disks.FindEfiFor(parent)
		if err != nil {
			return fmt.Errorf("can not find EFI. %v", err)
		}

		fmt.Printf("Found EFI: %v\n", efi.DeviceIdentifier)

		if efi.IsMounted() {
			return fmt.Errorf("EFI %s is already mounted", info.DeviceIdentifier)
		}

		osVersion, err := disk.GetOsVersion()
		if err != nil {
			return fmt.Errorf("can not get MacOS version. %v", err)
		}

		fmt.Printf("MacOS version: %s\n", osVersion.String())

		out, err := disk.MountDisk(efi, osVersion.GreaterThan(sudoMountOsVersion))
		if err != nil {
			return fmt.Errorf("can not mount disk %s", efi.DeviceIdentifier)
		}

		fmt.Printf("disk mounted: %s\n", out)
		return nil
	},
}
