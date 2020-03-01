// +build linux

package vmdetect

import (
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
		"kvm",
		"qemu",
		"vmware",
		"vmw",
		"oracle",
		"xen",
		"bochs",
		"parallels",
		"bhyve",
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

	return DoesFileContain(file, "Hypervisor detected")
}

/*
Checks if UML is being used
https://en.wikipedia.org/wiki/User-mode_Linux
*/
func checkUML() bool {

	file, err := os.Open("/proc/cpuinfo")

	if err != nil {
		PrintError(err)
		return false
	}

	defer file.Close()

	return DoesFileContain(file, "User Mode Linux")
}

/*
Some GNU/Linux distributions expose /proc/sysinfo containing potential VM info
https://www.ibm.com/support/knowledgecenter/en/linuxonibm/com.ibm.linux.z.lhdd/lhdd_t_sysinfo.html
*/
func checkSysInfo() bool {
	file, err := os.Open("/proc/sysinfo")

	if err != nil {
		PrintError(err)
		return false
	}

	defer file.Close()

	return DoesFileContain(file, "VM00")
}

/*
Public function returning true if a VM is detected.
If so, a non-empty string is also returned to tell how it was detected.
*/
func IsRunningInVirtualMachine() (bool, string) {

	if cpuid.CPU.VM() {
		return true, "CPU Vendor (assembly instructions)"
	}

	if checkUML() {
		return true, "CPU Vendor (/proc/cpuinfo)"
	}

	if checkSysInfo() {
		return true, "System Information (/proc/sysinfo)"
	}

	if checkDMITable() {
		return true, "DMI Table (/sys/class/dmi/id/*)"
	}

	if checkKernelRingBuffer() {
		return true, "Kernel Ring Buffer (/dev/kmsg)"
	}

	return false, "nothing"
}
