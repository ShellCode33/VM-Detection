// +build linux

package vmdetect

import (
	"bytes"
	"os"
	"os/exec"
)

/*
	Checks if the DMI table contains vendor strings of known VMs.
*/
func checkDMITable() bool {
	// TODO : instead of running a command, read files in /sys/class/dmi/id/* and look for vendor strings below
	output, err := exec.Command("dmidecode").Output()

	if err != nil {
		PrintError(err)
		return false
	}

	return bytes.Contains(output, []byte("innotek")) ||
		bytes.Contains(output, []byte("VirtualBox")) ||
		bytes.Contains(output, []byte("vbox"))
}

/*
	Checks printk messages to see if Linux detected an hypervisor.
*/
func checkKernelRingBuffer() bool {

	file, err := os.Open("/dev/kmsg")

	if err != nil {
		PrintError(err)
		return false
	}

	buffer := make([]byte, 1024*8)

	// Only reads the first 100 lines (reading character device in Go is annoying)
	for i := 0; i < 100; i++ {
		if _, err = file.Read(buffer); err != nil {
			PrintError(err)
			return false
		}

		if bytes.Contains(buffer, []byte("Hypervisor detected")) {
			return true
		}
	}

	return false
}

/*
	Public function returning true if a VM is detected.
	If so, a non-empty string is also returned to tell how it was detected.
*/
func IsRunningInVirtualMachine() (bool, string) {

	if checkDMITable() {
		return true, "DMI Table"
	}

	if checkKernelRingBuffer() {
		return true, "Kernel Ring Buffer"
	}

	return false, "nothing"
}
