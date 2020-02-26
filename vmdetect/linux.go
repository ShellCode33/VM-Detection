// +build linux darwin

package vmdetect

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
)

/*
	Checks if the program is being run using root.
 */
func isRunningWithAdminRights() (bool, error) {
	if currentUser, err := user.Current(); err != nil {
		return false, err
	} else {
		return currentUser.Uid == "0", nil
	}
}

/*
	Tries to vmdetect VM using privileged access.
 */
func privilegedChecks() (bool, string, error) {

	output, err := exec.Command("dmidecode").Output()

	if err == nil &&
		(bytes.Contains(output, []byte("innotek")) ||
		bytes.Contains(output, []byte("VirtualBox")) ||
		bytes.Contains(output, []byte("vbox"))){
		return true, "dmidecode", nil
	}

	output, err = exec.Command("dmesg").Output()

	if err == nil && bytes.Contains(output, []byte("Hypervisor detected")) {
		return true, "dmesg", nil
	}

	return false, "", nil
}

/*
	Tries to vmdetect VM using unprivileged access.
 */
func unprivilegedChecks() (bool, string, error) {
	output, err := exec.Command("hostnamectl").Output()

	if err == nil && bytes.Contains(output, []byte(" vm\n")) {
		return true, "hostnamectl", nil
	}

	return false, "", nil
}

/*
	Public function returning true if a VM is detected.
	If so, a non-empty string is also returned to tell how it was detected.
 */
func IsRunningInVirtualMachine() (bool, string, error) {
	isAdmin, err := isRunningWithAdminRights()

	if err != nil {
		return false, "", err
	}

	if isAdmin {
		vmDetected, reason, err := privilegedChecks()

		if err != nil {
			return false, "", err
		}

		if vmDetected {
			return true, reason, nil
		}
	} else {
		fmt.Println("[WARNING] Running as unprivileged user")
	}

	return unprivilegedChecks()
}