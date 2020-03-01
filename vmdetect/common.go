package vmdetect

import (
	"bufio"
	"fmt"
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

func DoesFileContain(file *os.File, stringToBeFound string) bool {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {

			if !os.IsTimeout(err) && err != io.EOF {
				PrintError(err)
			}

			return false
		}

		if strings.Contains(line, stringToBeFound) {
			return true
		}
	}
}
