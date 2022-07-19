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

	"github.com/sebnyberg/net/packet"
)

var (
	ifaceBufSize = 0
)

// Network is a set of interfaces (bound together by links). Currently, its
// role is to dynamically allocate IP addresses for interfaces.
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

// Node represents a node in one or more networks. Depending on how it manages
// traffic, it may be a router, switch, or perhaps a PC.
//
// When a packet arrives on an attached interface,
type Node struct {
	// Name contains the name of the node in the network
	Name string

	Interfaces []*Interface

	HandleIngress       PacketHandler
	HandleLocalDelivery PacketHandler
	HandleEgress        PacketHandler
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
	n.Interfaces = append(n.Interfaces, &netif)
	if net.interfaces == nil {
		net.interfaces = make(map[string]*Interface, 1)
	}
	net.interfaces[ifname] = &netif

	// Todo use node-level message handler instead of this dummy print
	go func() {
		for msg := range netif.Recv() {
			// Decode packet
			pkt, err := packet.Decode(msg)
			if err != nil {
				log.Println("failed to decode packet")
				continue
			}
			np := &NodePacket{
				Packet:   pkt,
				SourceIF: &netif,
			}
			n.handleIngress(np)
		}
	}()
	return &netif
}

// handleIngress goes through ingress handlers, reacting to their verdict.
// If the verdict is Accept, routing will either be routed to the destination
// interface (if set), or locally (if destination IF is unset).
func (n *Node) handleIngress(pkt *NodePacket) {
	// Todo: support chains of packet handlers
	for {
		ver := n.HandleIngress(pkt)
		switch ver {
		case VerdictAccept:
			goto accept
		case VerdictDrop:
			goto drop
		case VerdictRepeat:
			// Keep going
		default:
			log.Fatalln("invalid ingress verdict", ver)
		}
	}

drop:
	log.Println("dropping node packet due to ingress verdict")
	return

accept:
	if pkt.DestIF == nil {
		n.handleLocalDelivery(pkt)
	}
	n.handleEgress(pkt)
}

func (n *Node) handleLocalDelivery(pkt *NodePacket) {
	log.Fatalln("handle local delivery not implemented")
}

func (n *Node) handleEgress(pkt *NodePacket) {
	log.Fatalln("handle egress not implemented")
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

func allocHW() net.HardwareAddr {
	var addr [6]byte
	n, err := rand.Read(addr[:])
	if err != nil || n != 6 {
		panic(err)
	}
	return net.HardwareAddr(addr[:])
}

type Verdict uint8

const (
	// VerdictAccept continues packet iteration
	VerdictAccept Verdict = 0

	// VerdictDrop drops the packet immediately.
	VerdictDrop Verdict = 1

	// VerdictRepeat restarts packet iteration.
	// It is useful for when the packet contents have changed in some way.
	VerdictRepeat Verdict = 2

	// Todo: add other netfilter verdicts
	// https://netfilter.org/projects/libnetfilter_queue/doxygen/html/group__Queue.html
)

// PacketHandler is called on ingress, local delivery, and egress.
type PacketHandler func(p *NodePacket) Verdict

// NodePacket describes a packet flowing through a node's routing system.
// It is roughly equivalent to sk_buffer in Linux.
type NodePacket struct {
	// Packet contains a decoded Packet.
	// Layers of the packet can be manipulated (e.g. in the case of NAT), but be
	// wary of modification order.
	Packet packet.Packet

	// SourceIF points to the source interface.
	SourceIF *Interface

	// DestIF points to the destination interface.
	// If DestIF is nil at the end of ingress, then HandleLocalDelivery() is
	// called.
	DestIF *Interface
}
