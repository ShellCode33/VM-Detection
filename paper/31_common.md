\newpage

## Cross-platform solutions

When developing a malware, people usually target an operating system, as it can
be pretty difficult to build something that works as expected on each
environment. Moreover, releasing variants of the same binary (compiled for the
different environments) can facilitate the work of malware analysts. Despite 
these considerations, we aimed our researches toward cross-platform solutions
in order to mutualize efforts.

### Networking

Network adapters usually come with a MAC address (*MAC* stands for Media Access
Control, referring to the lowest part of the OSI model) which can be used to
identify its vendor. The first half of the address (the first 3 bytes) are
booked by constructors with the IEEE (Institute of Electrical and Electronics
Engineers, an international organisation dedicated to the writing of standards
for new technologies) to make the OUI, an unique vendor identifier.

Most hypervisors have an OUI so that it makes the network adapter easily
recognizable for the guest system. So if the a system sees such an OUI on its
network adapter, it is highly likely that it is a virtualized guest.

### Using CPUID

THe `CPUID` instruction has been introduced with Intel's *x86* architecture to
allow CPU discovery by the operating system. This way, the system can adapt its
behaviour to the characteristics of the processor. The use of this instruction
has been extended in 2008 to allow the hypervisor to "interact" with the guest
and thus optimizing its performance. By watching specific values of certain 
registers &ndash; mostly EBX, ECX, EDX &ndash; we can deduce the hypervisor if 
any.

### Measuring resources availability

Finally, low resources may be an indication that the operating system is
running inside a sandbox or virtual machine. It surely cannot be used as the
only clue but it can lead you to investigate: most sandboxes are ran on the
laptop of the analyst, who often will give the fewest resources they can. This
is why we consider machines with low resources (below 3GB of RAM and 3 CPUs) to
be virtual machines.
