\newpage

## Windows

### Crawling the Registry Hive

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

### Looking for known files

Following a similar principle to the registries analysis, we can take a look to
the file system, searching for particular files. When installed, the 
aforementioned *guest addons* add some files to the disk. These files are known
and pretty easy to guess. A lot of malwares will would take a look at these.
Here is a quick example:

```golang
vmwarePath := []string{
    `c:\windows\system32\drivers\vmmouse.sys`,
    `c:\windows\system32\drivers\vmnet.sys`,
    `c:\windows\system32\drivers\vmxnet.sys`,
    `c:\windows\system32\drivers\vmhgfs.sys`,
    `c:\windows\system32\drivers\vmx86.sys`,
    `c:\windows\system32\drivers\hgfs.sys`
}
```