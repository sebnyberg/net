package netlab

// Package netlab is loosely based on Tailscale's natlab package. It is used to
// simulate a network to debug different network topologies and algorithms.

import (
	"net/netip"
)

// Packet is an UDP packet flowing through the network.
type Packet struct {
	Src, Dst netip.AddrPort
	Payload  []byte
}

// Network simulates a network
type Network struct {
	Name   string
	Prefix netip.Prefix

	machines map[netip.Addr]*Interface
}

// Interface is a network interface
type Interface struct {
	machine *Machine
	net     *Network
	name    string
	ips     []netip.Addr
}

func (f *Interface) Machine() *Machine {
	return f.machine
}

func (f *Interface) Network() *Network {
	return f.net
}

// Machine is a machine in the network
type Machine struct {
	Name string

	interfaces []*Interface
}
