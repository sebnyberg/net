package netlab

// Package netlab provides a basic laboration network.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/netip"
)

var (
	ifaceBufSize = 0
)

// Network is a set of interfaces (bound together by links). Currently, it's
// role is to dynamically allocate IP addresses for interfaces. You could think
// of it as ICANN.
type Network struct {
	Name   string
	Prefix netip.Prefix
	lastV4 netip.Addr

	interfaces map[string]*Interface
}

// allocIP allocates an IP address for a new NIC on the network.
func (n *Network) allocIP() netip.Addr {
	if n.Prefix == (netip.Prefix{}) {
		panic("nil prefix")
	}
	if !n.Prefix.IsValid() {
		panic("invalid network prefix")
	}
	if n.lastV4 == (netip.Addr{}) {
		n.lastV4 = n.Prefix.Addr()
	}
	ip := n.lastV4
	n.lastV4 = n.lastV4.Next()
	if !n.Prefix.Contains(n.lastV4) {
		panic("no more ip addresses")
	}
	return ip
}

// Node is a node in the network. It could be a machine or perhaps a switch.
type Node struct {
	Name string

	interfaces []*Interface
}

// Attach attaches an interface to a node, granting it an IP (and MAC?) address.
// To link two interfaces together, use if.Link(other).
func (n *Node) Attach(ifname string, net *Network) *Interface {
	netif := Interface{
		Name: ifname,
		node: n,
		net:  net,
		link: nil,
		recv: make(chan []byte, ifaceBufSize),
		ip:   net.allocIP(),
		mac:  allocHW(),
	}
	n.interfaces = append(n.interfaces, &netif)
	if net.interfaces == nil {
		net.interfaces = make(map[string]*Interface, 1)
	}
	net.interfaces[ifname] = &netif

	// Todo use node-level message handler instead of this dummy print
	go func() {
		for msg := range netif.Recv() {
			log.Printf("received message of size %v on interface %v", len(msg), netif.Name)
		}
	}()
	return &netif
}

// Interface is a network interface that puts a machine onto a network via a
// link. To attach an interface to a node, use node.Attach().
// To link two interfaces together, use if.Link(other).
type Interface struct {
	Name string

	ip   netip.Addr
	mac  net.HardwareAddr
	node *Node
	net  *Network
	link *Link
	recv chan []byte
}

func (f *Interface) IP() netip.Addr {
	if f == nil {
		panic("nil interface")
	}
	return f.ip
}

func (f *Interface) IsUp() bool { return f.link != nil }

func (f *Interface) Link(toif *Interface) error {
	if f == nil || toif == nil {
		return errors.New("nil interface")
	}
	if f.net == nil || toif.net == nil {
		return errors.New("nil interface network")
	}
	if f.link != nil || toif.link != nil {
		return errors.New("link already exists")
	}
	if f.net.Name != toif.net.Name {
		return errors.New("must link within the same network")
	}
	link := &Link{
		if1: f,
		if2: toif,
	}
	f.link = link
	toif.link = link
	return nil
}

// Send sends a block of bytes to the interface.
func (f *Interface) Send(ctx context.Context, buf []byte) {
	if f.link == nil {
		panic("nil link")
	}
	other := f.link.if1
	if other == f {
		other = f.link.if2
	}
	select {
	case <-ctx.Done():
		log.Println("send timed out")
	case other.recv <- buf:
		fmt.Println("sent message")
	}
}

// Recv receives a block of bytes from the interface.
func (f *Interface) Recv() chan []byte {
	other := f.link.if1
	if other == f {
		other = f.link.if2
	}
	return other.recv
}

// Link binds two interfaces together.
type Link struct {
	if1, if2 *Interface
}

// Packet is UDP packet flowing through the network.
type Packet struct {
	src, dst netip.AddrPort

	layerTypes int    // bitmap of parsed layer types
	payload    []byte // entire packet
}

func allocHW() net.HardwareAddr {
	var addr [6]byte
	n, err := rand.Read(addr[:])
	if err != nil || n != 6 {
		panic(err)
	}
	return net.HardwareAddr(addr[:])
}
