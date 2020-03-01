# VM Detection

This project is a Go implementation of well-known techniques trying to detect if the program is being run in a virtual machine.

Why doing this in Go ? Because there are many C programs already doing this, but none written in pure Go.

See the [paper](https://github.com/ShellCode33/VM-Detection/blob/master/paper/paper.pdf) for more details.

## Usage

First download the package
```bash
$ go get github.com/ShellCode33/VM-Detection/vmdetect
```

Then see [main.go](https://github.com/ShellCode33/VM-Detection/blob/master/main.go) to use it in your own project.

To build the paper, be sure to have Docker installed and run the following command inside the paper directory:

```bash
$ docker run 
```

## GNU/Linux techniques

- Look for CPU vendor by trying out different assembly instructions ([cpuid](https://github.com/klauspost/cpuid/))
- Look for known strings in the DMI table (`/sys/class/dmi/id/*`)
- Look for hints in the kernel ring buffer (`/dev/kmsg`)

## Windows techniques

Coming soon...

## Resources

[systemd-detect-virt source code](https://github.com/systemd/systemd/blob/master/src/basic/virt.c)

[Malware evasion techniques](https://www.deepinstinct.com/2019/10/29/malware-evasion-techniques-part-2-anti-vm-blog/)
