package vmdetect

import "fmt"

func PrintError(loggee interface{}) {
	fmt.Printf("[x] %v\n", loggee)
}

func PrintWarning(loggee interface{}) {
	fmt.Printf("[!] %v\n", loggee)
}
