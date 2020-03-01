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

func DoesFileContain(file *os.File, stringToBeFound string) bool {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {

			if err == io.EOF {
				PrintError(file.Name() + " didn't match")
			} else if !os.IsTimeout(err) {
				PrintError(err)
			}

			return false
		}

		if strings.Contains(line, stringToBeFound) {
			return true
		}
	}
}
