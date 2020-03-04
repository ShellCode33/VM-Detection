# VM Detection

This project is a Go implementation of well-known techniques trying to detect if the program is being run in a virtual machine.
There are many C programs already doing this, but none written in pure Go.

See this [paper](https://github.com/ShellCode33/VM-Detection/blob/master/paper/paper.pdf) for more details.

## Usage

First download the package
```bash
$ go get github.com/ShellCode33/VM-Detection/vmdetect
```

Then see [main.go](https://github.com/ShellCode33/VM-Detection/blob/master/main.go) to use it in your own project.

This project is compatible for both Linux and Windows, you can use the following command to cross-compile it :
```bash
$ GOOS=windows go build main.go
$ file main.exe
```

## Common techniques

- Look for known mac address prefix
- Look for known interface names
- Look at CPU features using cpuid instruction ([cpuid](https://github.com/klauspost/cpuid/))

## GNU/Linux techniques

- Look for known strings in the DMI table `/sys/class/dmi/id/*`
- Look for hints in the kernel ring buffer `/dev/kmsg`
- Look for known LKM - Loadable Kernel Modules - `/proc/modules`
- Check existence of known files

## Windows techniques

- Check existence of known registry keys
- Look for known strings in some registry key's content
- Check existence of known files

## Credits

Thanks to [@hippwn](https://twitter.com/hippwn) for its contribution

Thanks systemd for being [that awesome](https://github.com/systemd/systemd/blob/master/src/basic/virt.c).

Thanks to CheckPoint's researchers for their [wonderful website](https://evasions.checkpoint.com/)