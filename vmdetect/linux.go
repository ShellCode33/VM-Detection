// +build linux

package vmdetect

import (
	"bufio"
	"bytes"
	"github.com/klauspost/cpuid"
	"io/ioutil"
	"os"
	"time"
)

/*
	Checks if the DMI table contains vendor strings of known VMs.
*/
func checkDMITable() bool {

	//  /!\ All lowercase /!\
	blacklistDMI := []string{
		"innotek",
		"virtualbox",
		"vbox",
	}

	dmiPath := "/sys/class/dmi/id/"
	dmiFiles, err := ioutil.ReadDir(dmiPath)

	if err != nil {
		PrintError(err)
		return false
	}

	for _, dmiEntry := range dmiFiles {
		if !dmiEntry.Mode().IsRegular() {
			continue
		}

		dmiContent, err := ioutil.ReadFile(dmiPath + dmiEntry.Name())

		if err != nil {
			PrintError(err)
			continue
		}

		for _, entry := range blacklistDMI {
			// Lowercase comparison to prevent false negatives
			if bytes.Contains(bytes.ToLower(dmiContent), []byte(entry)) {
				return true
			}
		}

	}

	return false
}

/*
	Checks printk messages to see if Linux detected an hypervisor.
	https://github.com/torvalds/linux/blob/31cc088a4f5d83481c6f5041bd6eb06115b974af/arch/x86/kernel/cpu/hypervisor.c#L79
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
			if !os.IsTimeout(err) {
				PrintError(err)
			}

			return false
		}

		// Lowercase comparison to prevent false negatives
		if bytes.Contains(bytes.ToLower(line), []byte("hypervisor detected")) {
			return true
		}
	}
}

/*
	Public function returning true if a VM is detected.
	If so, a non-empty string is also returned to tell how it was detected.
*/
func IsRunningInVirtualMachine() (bool, string) {

	if cpuid.CPU.VM() {
		return true, "CPU Vendor"
	}

	if checkDMITable() {
		return true, "DMI Table"
	}

	if checkKernelRingBuffer() {
		return true, "Kernel Ring Buffer"
	}

	return false, "nothing"
}
