This project is an attempt to understand using XDP to redirect packets at the veth level at different points.

The initial test is to see if we can redirect traffic coming from one container directly to another container, by attaching an XDP program to the exit veth of the sending container (foo), and using a redirect map with a single entry of the receiving container (bar).

The environment is setup by running `scripts/setup.sh`, which will create two namespaces, assign IPs and add static ARP entries.

After building and running the `xdp-redirect-map` binary, I can confirm that the redirect map is populated, and that the entry is what I expect:

```
$ sudo bpftool map dump name intf_map
key: 00 00 00 00  value: 02 01 00 00
Found 1 element
$ ip link show bar
258: bar@if2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default qlen 1000
    link/ether 42:43:b3:86:5f:12 brd ff:ff:ff:ff:ff:ff link-netnsid 7

```
`0x0102` is decimal 258 -- so I wonder if I am running into an endianess issue I'm not understanding?

Unfortunately, I cannot seem to make it work yet.  Running `ip netns exec foons iperf -c 192.0.2.1 -u -b 1m -t 240`, I can then run `ip netns exec foons tcpdump -nei eth0` and see traffic going into the netns veth device, but I never see traffic arriving if I do `ip netns exec barns tcpdump -nei eth0`.

I know the traffic is being redirected, because `tcpdump -nei foo` on the main namespace doesn't show any traffic while the XDP is loaded.
