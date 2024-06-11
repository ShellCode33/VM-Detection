package main

import (
	"fmt"
	"github.com/ShellCode33/VM-Detection/vmdetect"
)

func main() {
	fmt.Println("Trying to detect if a VM is running...")

	isInsideVM, reason := vmdetect.IsRunningInVirtualMachine()

	if isInsideVM {
		fmt.Printf("\nVM detected thanks to %v\n", reason)
	} else {
		fmt.Println("\nNo VM detected")
	}

}
