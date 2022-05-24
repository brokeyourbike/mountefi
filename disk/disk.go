package disk

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	version "github.com/hashicorp/go-version"
	"howett.net/plist"
)

const efi = "efi"
const appleApfs = "apple_apfs"
const appleCoreStorage = "apple_corestorage"

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
	APFSContainerReference    string `plist:"APFSContainerReference"`
	APFSPhysicalStores        []struct {
		APFSPhysicalStore string `plist:"APFSPhysicalStore"`
	} `plist:"APFSPhysicalStores"`
}

func (d *DiskInfo) IsMounted() bool {
	return d.MountPoint != ""
}

func (d *DiskInfo) IsEfi() bool {
	return strings.ToLower(d.Content) == efi
}

func (d *DiskInfo) IsApfs() bool {
	return d.APFSContainerReference != ""
}

func (d *DiskInfo) IsApfsContainer() bool {
	return strings.ToLower(d.Content) == appleApfs
}

func (d *DiskInfo) IsCoreStorageContainer() bool {
	return strings.ToLower(d.Content) == appleCoreStorage
}

func (d *DiskInfo) GetApfsPhysicalStores() []string {
	stores := make([]string, 0)

	for _, s := range d.APFSPhysicalStores {
		stores = append(stores, s.APFSPhysicalStore)
	}

	return stores
}

type Disks struct {
	target []string
	list   map[string]DiskInfo
	mu     sync.RWMutex
}

func NewDisks(target []string) *Disks {
	d := &Disks{
		target: target,
		list:   make(map[string]DiskInfo, 0),
		mu:     sync.RWMutex{},
	}

	return d
}

func (d *Disks) Update() {
	var wg = sync.WaitGroup{}

	for i := range d.target {
		wg.Add(1)
		v := d.target[i]
		go func() {
			defer wg.Done()

			info, err := GetDiskInfo(v)
			if err != nil {
				fmt.Printf("err: %v", err)
				return
			}

			d.mu.Lock()
			defer d.mu.Unlock()

			d.list[info.DeviceIdentifier] = info
		}()
	}

	wg.Wait()
}

func (d *Disks) GetIdentifiers() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	ids := make([]string, 0)
	for i := range d.list {
		ids = append(ids, i)
	}
	return ids
}

// FindParentFor will try to the root physical storage.
func (d *Disks) FindParentFor(disk DiskInfo) (DiskInfo, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	possibleParents := []string{disk.ParentWholeDisk}

	if disk.IsApfs() {
		for _, s := range disk.GetApfsPhysicalStores() {
			// find the physical store in the list
			v, ok := d.list[s]
			if !ok {
				continue
			}

			// get it's parent
			possibleParents = append(possibleParents, v.ParentWholeDisk)
		}
	}

	for _, pp := range possibleParents {
		parent, ok := d.list[pp]
		if !ok {
			continue
		}

		if parent.IsApfs() {
			continue
		}

		if parent.WholeDisk {
			return parent, nil
		}
	}

	return DiskInfo{}, fmt.Errorf("can not found parent for %s", disk.DeviceIdentifier)
}

// FindEfiFor will try to find the EFI volume, owned by the provided disk.
func (d *Disks) FindEfiFor(parent DiskInfo) (DiskInfo, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, v := range d.list {
		if v.IsEfi() && parent.DeviceIdentifier == v.ParentWholeDisk {
			return v, nil
		}
	}

	return DiskInfo{}, fmt.Errorf("cannot found EFI for %s", parent.DeviceIdentifier)
}

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

	v := strings.TrimRight(string(raw), "\n")
	return version.NewSemver(v)
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

func MountDisk(disk DiskInfo, withSudo bool) (string, error) {
	cmd := exec.Command("diskutil", "mount", disk.DeviceIdentifier)
	if withSudo {
		cmd = exec.Command("sudo", "diskutil", "mount", disk.DeviceIdentifier)
	}
	raw, err := cmd.Output()
	return string(raw), err
}

func UnmountDisk(disk DiskInfo) (string, error) {
	raw, err := exec.Command("diskutil", "unmount", disk.DeviceIdentifier).Output()
	return string(raw), err
}
