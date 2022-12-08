#!/bin/bash
set +eux

sudo ip netns add foons
sudo ip netns add barns
sudo ip link add foo type veth peer eth0 netns foons
sudo ip link add bar type veth peer eth0 netns barns
sudo ip link set foo up
sudo ip link set bar up
sudo ip netns exec foons ip link set lo up
sudo ip netns exec foons ip link set dev eth0 address 76:E3:1B:BC:E1:E3
sudo ip netns exec foons ip link set eth0 up
sudo ip netns exec foons ip addr add 192.0.2.0/31 dev eth0
sudo ip netns exec foons ip neigh add 192.0.2.1 lladdr 76:E3:1B:BC:E1:E4 dev eth0
sudo ip netns exec barns ip link set lo up
sudo ip netns exec barns ip link set dev eth0 address 76:E3:1B:BC:E1:E4
sudo ip netns exec barns ip link set eth0 up
sudo ip netns exec barns ip addr add 192.0.2.1/31 dev eth0
sudo ip netns exec foons ip neigh add 192.0.2.0 lladdr 76:E3:1B:BC:E1:E3 dev eth0
