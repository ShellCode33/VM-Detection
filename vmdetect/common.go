package vmdetect

import (
	"bufio"
	"fmt"
	"github.com/klauspost/cpuid"
	"io"
	"os"
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
Tries to detect VMs using cross-platform techniques.
*/
func CommonChecks() (bool, string) {
	// https://lwn.net/Articles/301888/
	if cpuid.CPU.VM() {
		return true, "CPU Vendor (cpuid space)"
	}

	return false, "nothing"
}
