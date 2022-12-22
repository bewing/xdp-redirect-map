package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cilium/ebpf/link"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go redirect redirect.c -- -O2 -Wall

func main() {
	bpfObjs := redirectObjects{}
	if err := loadRedirectObjects(&bpfObjs, nil); err != nil {
		log.Fatalf("Could not load BPF objects: %s", err)
	}
	iface, err := net.InterfaceByName("foo")
	if err != nil {
		log.Fatalf("Could not find interface foo: %s", err)
	}
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   bpfObjs.XdpRedirectFunc,
		Interface: iface.Index,
		Flags:     link.XDPGenericMode,
	})
	defer l.Close()
	iface, err = net.InterfaceByName("bar")
	if err != nil {
		log.Fatalf("Could not find interface foo: %s", err)
	}
	if err := bpfObjs.IntfMap.Put(uint32(0), uint32(iface.Index)); err != nil {
		log.Fatalf("Could not insert into map: %s", err)
	}
	fmt.Println("XDP loaded and map updated. Press ctrl-c to exit")
	for {
		time.Sleep(time.Second)
	}
}
