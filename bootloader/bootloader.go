package bootloader

import (
	"os/exec"
)

func FindBootloaderID() (string, error) {
	return findOpencoreID()
}

func findOpencoreID() (string, error) {
	raw, err := exec.Command("nvram", "4D1FDA02-38C7-4A6A-9CC6-4BCCA8B30102:boot-path").Output()
	return string(raw), err
}
