package disk

import (
	"fmt"
	"os/exec"
	"strings"

	version "github.com/hashicorp/go-version"
	"howett.net/plist"
)

const efi = "efi"

// const appleApfs = "apple_apfs"
// const appleCorestorage = "apple_corestorage"

type DiskInfo struct {
	Content                   string `plist:"Content"`
	DeviceIdentifier          string `plist:"DeviceIdentifier"`
	DeviceNode                string `plist:"DeviceNode"`
	DiskUUID                  string `plist:"DiskUUID"`
	VolumeName                string `plist:"VolumeName"`
	VolumeUUID                string `plist:"VolumeUUID"`
	ParentWholeDisk           string `plist:"ParentWholeDisk"`
	WholeDisk                 bool   `plist:"WholeDisk"`
	MountPoint                string `plist:"MountPoint"`
	FilesystemName            string `plist:"FilesystemName"`
	FilesystemUserVisibleName string `plist:"FilesystemUserVisibleName"`
	FilesystemType            string `plist:"FilesystemType"`
}

func (d *DiskInfo) IsMounted() bool {
	return d.MountPoint != ""
}

func (d *DiskInfo) IsEfi() bool {
	return strings.ToLower(d.Content) == efi
}

type Disks []DiskInfo

type DiskList struct {
	AllDisks         []string `plist:"AllDisks"`
	VolumesFromDisks []string `plist:"VolumesFromDisks"`
	WholeDisks       []string `plist:"WholeDisks"`
}

// GetDiskutilPath executes `which diskutil` and returns path
func GetDiskutilPath() (string, error) {
	path, err := exec.Command("which", "diskutil").Output()
	return string(path), err
}

// GetOsVersion executes `sw_vers -productVersion` command and returns
// semver compliant version on success
func GetOsVersion() (*version.Version, error) {
	raw, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return nil, err
	}
	return version.NewSemver(string(raw))
}

// GetDiskList executes `diskutil list -plist` and returns DiskList
func GetDiskList() (list DiskList, err error) {
	raw, err := exec.Command("diskutil", "list", "-plist").Output()
	if err != nil {
		return
	}

	_, err = plist.Unmarshal(raw, &list)
	return
}

// GetDiskInfo executes `diskutil info -plist <disk>` command and returns DiskInfo
func GetDiskInfo(disk string) (diskInfo DiskInfo, err error) {
	raw, err := exec.Command("diskutil", "info", "-plist", disk).Output()
	if err != nil {
		return
	}

	_, err = plist.Unmarshal(raw, &diskInfo)
	return
}

func GetVolumeInfo(volume string) (diskInfo DiskInfo, err error) {
	if !strings.HasPrefix(volume, "/Volumes/") {
		volume = fmt.Sprintf("/Volumes/%s", volume)
	}
	return GetDiskInfo(volume)
}

func GetMountedVolumes() ([]string, error) {
	raw, err := exec.Command("ls", "-1", "/Volumes").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(raw), "\n"), "\n"), nil
}

func MountDisk(disk DiskInfo, withSudo bool) error {
	cmd := exec.Command("diskutil", "mount", disk.DeviceIdentifier)
	if withSudo {
		cmd = exec.Command("sudo", "diskutil", "mount", disk.DeviceIdentifier)
	}
	_, err := cmd.Output()
	return err
}

func UnmountDisk(disk DiskInfo) error {
	_, err := exec.Command("diskutil", "unmount", disk.DeviceIdentifier).Output()
	return err
}
