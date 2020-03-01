\newpage

# State of the art

In this part, we will dive into the concept of evasion. Malware authors are
ahead of the analyst and it is a lead they have to maintain. To this end, they
deploy more and more techniques to make harder the analysis and comprehension 
of their code. Against static analysis they use obfuscation, encryption and 
such, as for dynamic ones they use evasion.

The concept of evasion refers to all the techniques used by a malware to hide
its behaviour according to its environment. For instance, if a malware detects
a sandbox (from an analyst or an antivirus), it will make low profile to so as
not to arouse suspicion. One of the most known example of this is the Red Pill
demonstration (reference to the legendary film *The Matrix*) presented by
Joanna Rutkowska in 2004 &ndash; just two years before she presents the Blue
Pill attack which is a type of *hyperjacking*.

Red Pill is a small piece of code written in C that checks the address of the 
*Interrupt Descriptor Table* (IDT). The address of this table has to be 
modified by the hypervisor to avoid memory conflicts. Therefore, there is a 
correlation between the address being superior to `0xD0` and the fact of being 
executed in a virtual machines. This technique is however less efficient on
today's systems as those filter the access to certain zones of the memory, such
as the DTI. It still works on *QEMU* though.

```c
int swallow_redpill () 
{
  unsigned char m[2+4], rpill[] = "\x0f\x01\x0d\x00\x00\x00\x00\xc3";
  *((unsigned*)&rpill[3]) = (unsigned)m;
  ((void(*)())&rpill)();
  return (m[5]>0xd0) ? 1 : 0;
}
```

## Linux techniques

### The DMI table

DMI stands for *Desktop Management Interface*. It is a standard developed in 
the 90' with de goal of uniforming the tracking of the components in a computer
and abstracting them from the softwares supposed to run them. Parsing this
table can reveal practical information on the hardware used by the operating
system and possibly detect the presence of names specific to virtualized
environment, such as *vbox*, *virtualbox*, *oracle*, *qemu*, *kvm* and so on.

### Linux kernel's hypervisor detection

Linux's kernel comes with an hypervisor detection feature that can be used to
identify a potential hypervisor below the operating system. Based on this, we
easily can listen for the kernel event to see if an hypervisor has been
detected by the kernel:

```c
static inline const struct hypervisor_x86 * __init
detect_hypervisor_vendor(void)
{
	const struct hypervisor_x86 *h = NULL, * const *p;
	uint32_t pri, max_pri = 0;

	for (p = hypervisors; p < hypervisors + ARRAY_SIZE(hypervisors); p++) {
		if (unlikely(nopv) && !(*p)->ignore_nopv)
			continue;

		pri = (*p)->detect();
		if (pri > max_pri) {
			max_pri = pri;
			h = *p;
		}
	}

	if (h)
        // this line prints the hypervisor in the `/dev/kmsg` file
		pr_info("Hypervisor detected: %s\n", h->name);

	return h;
}
```

### Checking Linux's pseudo-filesystems

Linux provides a lot of information via a certain type of files (mostly in
`/proc`) that are generated at boot and modified during runtime. A lot of 
binaries use this directory like `ps`, `uname`, `lspci` and so on. These
information are really helpful when trying to identify wether or not you are
in a virtualized environment, like UML for instance. UML refers to the
aforementioned way of executing a Linux kernel in user-space. This can easily
be verified by looking for the string "User Mode Linux" in the file
`/proc/cpuinfo` which describes the CPU of the machine.

In the same way, a lot of these virtual *files* can provide information on the
environment, including &ndash; but not limited to &ndash; `/proc/sysinfo` (in
which some distribution expose data about virtual machines),
`/proc/device-tree` (that lists the devices on the machine), `/proc/xen` (a 
file created by the *Xen Server*) or `/proc/modules` (that contains information
about the loaded kernel modules, modules that are used by hypervisors to 
optimize the guests).

Like *procfs* (mounted in `/proc`), *sysfs* can be useful. Its role is to
provide to the user an access to the devices and their drivers. The file
`/sys/hypervisor/type`, for instance, is sometimes used to store information
about the hypervisor Linux is running on.


## Windows

<!-- TODO -->