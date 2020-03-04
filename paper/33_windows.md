\newpage

## Windows

On Windows, most configuration can be done through the *Registry Hive* &ndash;
some kind of database that contains every configuration option about either the
operating system itself, or any software that would like to store information
in it. A lot of indicators of hypervisors can be stored there, especially if
the *guest addons* (small pieces of software that are installed on the guest 
to allow interoperability between the guest and the host, permitting shared
clipboard, *drag'n'drop* and so on) are installed.

Most keys will be installed inside the `HKEY_LOCAL_MACHINE` register which
mostly contains information about hardware, security and such. Parsing its
content looking for particular patterns is efficient enough and quite a good 
indicator of the presence of an hypervisor if any. Here is an example of keys 
that we are looking for:

```golang
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
```