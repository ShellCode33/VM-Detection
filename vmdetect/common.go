package vmdetect

import (
	"bufio"
	"fmt"
	"github.com/klauspost/cpuid"
	"github.com/shirou/gopsutil/mem"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
)

func PrintError(loggee interface{}) {
	fmt.Printf("[x] %v\n", loggee)
}

func PrintWarning(loggee interface{}) {
	fmt.Printf("[!] %v\n", loggee)
}

func DoesFileExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		PrintError(err)
	}

	return !os.IsNotExist(err)
}

func DoesFileContain(file *os.File, stringsToBeFound ...string) bool {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {

			if !os.IsTimeout(err) && err != io.EOF {
				PrintError(err)
			}

			return false
		}

		for _, stringToBeFound := range stringsToBeFound {
			if strings.Contains(line, stringToBeFound) {
				return true
			}
		}
	}
}

/*
Tries to detect a VM using its network configuration.
*/
func checkNetworking() (bool, string) {

	blacklistedMacAddressPrefixes := []string{
		"00:1C:42", // Parallels
		"08:00:27", // VirtualBox
		"00:05:69", // |
		"00:0C:29", // | > VMWare
		"00:1C:14", // |
		"00:50:56", // |
		"00:16:E3", // Xen
	}

	interfaces, err := net.Interfaces()

	if err != nil {
		return false, err.Error()
	}

	for _, iface := range interfaces {

		macAddr := iface.HardwareAddr.String()
		if macAddr != "" {
			for _, prefix := range blacklistedMacAddressPrefixes {
				if strings.HasPrefix(macAddr, prefix) {
					return true, fmt.Sprintf("Known MAC address prefix (%v)", prefix)
				}
			}
		}

		if iface.Name == "Vmware" {
			return true, "Vmware found as network interface name"
		}
	}

	return false, ""
}

/*
Tries to detect VMs using cross-platform techniques.
*/
func CommonChecks() (bool, string) {

	// https://lwn.net/Articles/301888/
	if cpuid.CPU.VM() {
		return true, "CPU Vendor (cpuid space)"
	}

	if vmDetected, how := checkNetworking(); vmDetected {
		return true, how
	}

	vmStat, err := mem.VirtualMemory()

	if err != nil {
		PrintError(err)
	} else if runtime.NumCPU() < 3 && vmStat.Total < 2048000 {
		return true, fmt.Sprintf("Low resources detected (%v CPU and %v bytes of RAM", runtime.NumCPU(), vmStat.Total)
	}

	return false, "nothing"
}
