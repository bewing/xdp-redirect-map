This project is an attempt to understand using XDP to redirect packets at the veth level at different points.

The initial test is to see if we can redirect traffic coming from one container directly to another container, by attaching an XDP program to the exit veth of the sending container (foo), and using a redirect map with a single entry of the receiving container (bar).

The environment is setup by running `scripts/setup.sh`, which will create two namespaces, assign IPs and add static ARP entries:

```
┌─────────────────┐   ┌─────────────────┐
│ foons           │   │ barns           │
│       ┌──────┐  │   │ ┌──────┐        │
│       │eth0  │  │   │ │eth0  │        │
└───────┴──┬───┴──┘   └─┴──┬───┴────────┘
           │               │
           │               │
           │               │
     ┌──┬──┴───┬────────┬──┴───┬──┐
     │  │foo   │        │bar   │  │
     │  └──────┘        └──────┘  │
     │                            │
     │            NS0             │
     │                            │
     └────────────────────────────┘
```

After building and running the `xdp-redirect-map` binary, I can confirm that the redirect map is populated, and that the entry is what I expect:

```
$ sudo bpftool map dump name intf_map
key: 00 00 00 00  value: 02 01 00 00
Found 1 element
$ ip link show bar
258: bar@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default qlen 1000
    link/ether 42:43:b3:86:5f:12 brd ff:ff:ff:ff:ff:ff link-netnsid 7

```

This code works as-is attaching to `foo` with `link.XDPGenericMode` (skb) mode.  If attempting to attach with `link.XDPDriverMode` (the default),
you **MUST** attach an XDP program of some sort to the `eth0` in the `barns` namespace, as noted in this paper:
https://www.files.netdevconf.info/f/a63b274e50f943a0a474
