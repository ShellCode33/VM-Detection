package main

import (
	"VMDetect/vmdetect"
	"fmt"
)

func main() {
	fmt.Println("Trying to vmdetect if a VM is running...")

	isInsideVM, reason, err := vmdetect.IsRunningInVirtualMachine()

	if err != nil {
		panic(err)
	}

	if isInsideVM {
		fmt.Printf("VM detected thanks to %v\n", reason)
	} else {
		fmt.Println("No VM has been detected")
	}

}
