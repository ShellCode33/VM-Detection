// +build windows

package vmdetect

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func extractKeyTypeFrom(registryKey string) (registry.Key, string, error) {
	firstSeparatorIndex := strings.Index(registryKey, string(os.PathSeparator))
	keyTypeStr := registryKey[:firstSeparatorIndex]
	keyPath := registryKey[firstSeparatorIndex+1:]

	var keyType registry.Key
	switch keyTypeStr {
	case "HKLM":
		keyType = registry.LOCAL_MACHINE
		break
	case "HKCR":
		keyType = registry.CLASSES_ROOT
		break
	case "HKCU":
		keyType = registry.CURRENT_USER
		break
	case "HKU":
		keyType = registry.USERS
		break
	case "HKCC":
		keyType = registry.CURRENT_CONFIG
		break
	default:
		return keyType, "", errors.New(fmt.Sprintf("Invalid keytype (%v)", keyTypeStr))
	}

	return keyType, keyPath, nil
}

func doesRegistryKeyContain(registryKey string, expectedSubString string) bool {

	keyType, keyPath, err := extractKeyTypeFrom(registryKey)

	if err != nil {
		PrintError(err)
		return false
	}

	keyPath, keyName := filepath.Split(keyPath)

	keyHandle, err := registry.OpenKey(keyType, keyPath, registry.QUERY_VALUE)

	if err != nil {
		PrintError(fmt.Sprintf("Cannot open %v : %v", registryKey, err))
		return false
	}

	defer keyHandle.Close()

	valueFound, _, err := keyHandle.GetStringValue(keyName)

	if err != nil {
		PrintError(err)
	}

	return strings.Contains(valueFound, expectedSubString)
}

func doesRegistryKeyExist(registryKey string) bool {

	subkeyPrefix := ""

	// Handle trailing wildcard
	if registryKey[len(registryKey)-1:] == "*" {
		registryKey, subkeyPrefix = path.Split(registryKey)
		subkeyPrefix = subkeyPrefix[:len(subkeyPrefix)-1] // remove *
	}

	keyType, keyPath, err := extractKeyTypeFrom(registryKey)

	if err != nil {
		PrintError(err)
		return false
	}

	keyHandle, err := registry.OpenKey(keyType, keyPath, registry.QUERY_VALUE)

	if err != nil {
		PrintError(fmt.Sprintf("Cannot open %v : %v", registryKey, err))
		return false
	}

	defer keyHandle.Close()

	// If a wildcard has been specified...
	if subkeyPrefix != "" {
		// ... we look for sub-keys to see if one exists
		subKeys, err := keyHandle.ReadSubKeyNames(0xFFFF)

		if err != nil {
			PrintError(err)
			return false
		}

		for _, subKeyName := range subKeys {
			if strings.HasPrefix(subKeyName, subkeyPrefix) {
				return true
			}
		}

		return false
	} else {
		// The registryKey we were looking for has been found
		return true
	}
}

func checkRegistry() (bool, string) {

	hyperVKeys := []string{
		`HKLM\SOFTWARE\Microsoft\Hyper-V`,
		`HKLM\SOFTWARE\Microsoft\VirtualMachine`,
		`HKLM\SOFTWARE\Microsoft\Virtual Machine\Guest\Parameters`,
		`HKLM\SYSTEM\ControlSet001\Services\vmicheartbeat`,
		`HKLM\SYSTEM\ControlSet001\Services\vmicvss`,
		`HKLM\SYSTEM\ControlSet001\Services\vmicshutdown`,
		`HKLM\SYSTEM\ControlSet001\Services\vmicexchange`,
	}

	parallelsKeys := []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_1AB8*`,
	}

	virtualBoxKeys := []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_80EE*`,
		`HKLM\HARDWARE\ACPI\DSDT\VBOX__`,
		`HKLM\HARDWARE\ACPI\FADT\VBOX__`,
		`HKLM\HARDWARE\ACPI\RSDT\VBOX__`,
		`HKLM\SOFTWARE\Oracle\VirtualBox Guest Additions`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxGuest`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxMouse`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxService`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxSF`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxVideo`,
	}

	virtualPCKeys := []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_5333*`,
		`HKLM\SYSTEM\ControlSet001\Services\vpcbus`,
		`HKLM\SYSTEM\ControlSet001\Services\vpc-s3`,
		`HKLM\SYSTEM\ControlSet001\Services\vpcuhub`,
		`HKLM\SYSTEM\ControlSet001\Services\msvmmouf`,
	}

	vmwareKeys := []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_15AD*`,
		`HKCU\SOFTWARE\VMware, Inc.\VMware Tools`,
		`HKLM\SOFTWARE\VMware, Inc.\VMware Tools`,
		`HKLM\SYSTEM\ControlSet001\Services\vmdebug`,
		`HKLM\SYSTEM\ControlSet001\Services\vmmouse`,
		`HKLM\SYSTEM\ControlSet001\Services\VMTools`,
		`HKLM\SYSTEM\ControlSet001\Services\VMMEMCTL`,
		`HKLM\SYSTEM\ControlSet001\Services\vmware`,
		`HKLM\SYSTEM\ControlSet001\Services\vmci`,
		`HKLM\SYSTEM\ControlSet001\Services\vmx86`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\CdRomNECVMWar_VMware_IDE_CD*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\CdRomNECVMWar_VMware_SATA_CD*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\DiskVMware_Virtual_IDE_Hard_Drive*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\DiskVMware_Virtual_SATA_Hard_Drive*`,
	}

	xenKeys := []string{
		`HKLM\HARDWARE\ACPI\DSDT\xen`,
		`HKLM\HARDWARE\ACPI\FADT\xen`,
		`HKLM\HARDWARE\ACPI\RSDT\xen`,
		`HKLM\SYSTEM\ControlSet001\Services\xenevtchn`,
		`HKLM\SYSTEM\ControlSet001\Services\xennet`,
		`HKLM\SYSTEM\ControlSet001\Services\xennet6`,
		`HKLM\SYSTEM\ControlSet001\Services\xensvc`,
		`HKLM\SYSTEM\ControlSet001\Services\xenvdb`,
	}

	// TODO : fill with https://evasions.checkpoint.com/techniques/registry.html#check-if-keys-contain-strings
	blacklistedValuesPerKeyPerVendor := map[string]map[string]string{
		"Anubis": {
			`HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\ProductID`:    "76487-337-8429955-22614",
			`HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProductID`: "76487-337-8429955-22614",
		},
		"QEMU": {
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: "QEMU",
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`:                                           "QEMU",
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:                                            "QEMU",
			`HKLM\HARDWARE\Description\System\BIOS\SystemManufacturer`:                                     "QEMU",
		},
		"VirtualBox": {
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: "VBOX",
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 1\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: "VBOX",
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 2\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: "VBOX",
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`:                                           "VBOX",
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:                                            "VIRTUALBOX",
			`HKLM\HARDWARE\Description\System\BIOS\SystemProductName`:                                      "VIRTUAL",
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\DeviceDesc`:                                      "VBOX",
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\FriendlyName`:                                    "VBOX",
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\DeviceDesc`:                                      "VBOX",
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\FriendlyName`:                                    "VBOX",
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\DeviceDesc`:                                      "VBOX",
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\FriendlyName`:                                    "VBOX",
			`HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\SystemProductName`:                    "VIRTUAL",
		},
	}

	allKeys := [][]string{hyperVKeys, parallelsKeys, virtualBoxKeys, virtualPCKeys, vmwareKeys, xenKeys}

	for _, keys := range allKeys {
		for _, key := range keys {
			if doesRegistryKeyExist(key) {
				return true, key
			}
		}
	}

	for /*vendor*/ _, registryValuesPerPath := range blacklistedValuesPerKeyPerVendor {
		for registryPath, expectedValue := range registryValuesPerPath {
			if doesRegistryKeyContain(registryPath, expectedValue) {
				return true, registryPath + " contains " + expectedValue
			}
		}
	}

	return false, "none"
}

func checkFileSystem() (bool, string) {
	// check for known path on the filesystem, either files or directories
	generalPath := []string{
		`c:\take_screenshot.ps1`,
		`c:\loaddll.exe`,
		`c:\symbols\aagmmc.pdb`,
	}

	prlPath := []string{
		`c:\windows\system32\drivers\prleth.sys`,
		`c:\windows\system32\drivers\prlfs.sys`,
		`c:\windows\system32\drivers\prlmouse.sys`,
		`c:\windows\system32\drivers\prlvideo.sys`,
		`c:\windows\system32\drivers\prltime.sys`,
		`c:\windows\system32\drivers\prl_pv32.sys`,
		`c:\windows\system32\drivers\prl_paravirt_32.sys`,
	}

	vboxPath := []string{
		`c:\windows\system32\drivers\VBoxMouse.sys`,
		`c:\windows\system32\drivers\VBoxGuest.sys`,
		`c:\windows\system32\drivers\VBoxSF.sys`,
		`c:\windows\system32\drivers\VBoxVideo.sys`,
		`c:\windows\system32\vboxdisp.dll`,
		`c:\windows\system32\vboxhook.dll`,
		`c:\windows\system32\vboxmrxnp.dll`,
		`c:\windows\system32\vboxogl.dll`,
		`c:\windows\system32\vboxoglarrayspu.dll`,
		`c:\windows\system32\vboxoglcrutil.dll`,
		`c:\windows\system32\vboxoglerrorspu.dll`,
		`c:\windows\system32\vboxoglfeedbackspu.dll`,
		`c:\windows\system32\vboxoglpackspu.dll`,
		`c:\windows\system32\vboxoglpassthroughspu.dll`,
		`c:\windows\system32\vboxservice.exe`,
		`c:\windows\system32\vboxtray.exe`,
		`c:\windows\system32\VBoxControl.exe`,
	}

	vmwarePath := []string{
		`c:\windows\system32\drivers\vmmouse.sys`,
		`c:\windows\system32\drivers\vmnet.sys`,
		`c:\windows\system32\drivers\vmxnet.sys`,
		`c:\windows\system32\drivers\vmhgfs.sys`,
		`c:\windows\system32\drivers\vmx86.sys`,
		`c:\windows\system32\drivers\hgfs.sys`,
	}

	virtualpcPath := []string{
		`c:\windows\system32\drivers\vmsrvc.sys`,
		`c:\windows\system32\drivers\vpc-s3.sys`,
	}

	allPath := [][]string{virtualpcPath, prlPath, vmwarePath, vboxPath, generalPath}

	for _, paths := range allPath {
		for _, p := range paths {
			if DoesFileExist(p) {
				return true, p
			}
		}
	}

	return false, "none"
}

/*
	Public function returning true if a VM is detected.
	If so, a non-empty string is also returned to tell how it was detected.
*/
func IsRunningInVirtualMachine() (bool, string) {
	if vmDetected, how := CommonChecks(); vmDetected {
		return vmDetected, how
	}

	if vmDetected, registryKey := checkRegistry(); vmDetected {
		return vmDetected, fmt.Sprintf("Registry key (%v)", registryKey)
	}

	if vmDetected, filePath := checkFileSystem(); vmDetected {
		return vmDetected, fmt.Sprintf("Known path found (%v)", filePath)
	}

	return false, "nothing"
}
