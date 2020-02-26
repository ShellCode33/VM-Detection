# VM Detection

This project is a Go implementation of well-known techniques trying to detect if the program is being run in a virtual machine.

Why doing this in Go ? Because there are many C programs already doing this, but none written in pure Go.

## GNU/Linux techniques

- Look for known strings in the DMI table (/dev/mem)
- Look for hints in the kernel ring buffer (/dev/kmsg)
- Look for virtual chassis in systemd configuration

## Windows techniques

Coming soon...

## Resources

![systemd-detect-virt source code](https://github.com/systemd/systemd/blob/master/src/basic/virt.c)
![Malware evasion techniques](https://www.deepinstinct.com/2019/10/29/malware-evasion-techniques-part-2-anti-vm-blog/)
