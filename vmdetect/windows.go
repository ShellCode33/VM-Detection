// +build windows

package vmdetect

/*
	Checks if the program is being run in a privileged context.
*/
func isRunningWithAdminRights() (bool, error) {
	return false, nil
}

/*
	Tries to vmdetect VM using privileged access.
*/
func privilegedChecks() (bool, string, error) {
	return false, "", nil
}

/*
	Tries to vmdetect VM using unprivileged access.
*/
func unprivilegedChecks() (bool, string, error) {
	return false, "", nil
}

/*
	Public function returning true if a VM is detected.
	If so, a non-empty string is also returned to tell how it was detected.
*/
func IsRunningInVirtualMachine() (bool, string, error) {
	isAdmin, err := isRunningWithAdminRights()

	if err != nil {
		return false, "", err
	}

	if isAdmin {
		vmDetected, reason, err := privilegedChecks()

		if err != nil {
			return false, "", err
		}

		if vmDetected {
			return true, reason, nil
		}
	} else {
		fmt.Println("[WARNING] Running as unprivileged user")
	}

	return unprivilegedChecks()
}