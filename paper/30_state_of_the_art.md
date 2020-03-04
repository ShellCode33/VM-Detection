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
