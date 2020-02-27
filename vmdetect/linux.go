// +build linux

package vmdetect

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"time"
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

	defer file.Close()

	// Set a read timeout because otherwise reading kmsg (which is a character device) will block
	if err = file.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		PrintError(err)
		return false
	}

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if ! os.IsTimeout(err) {
				PrintError(err)
			}

			return false
		}

		if bytes.Contains(line, []byte("Hypervisor detected")) {
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
