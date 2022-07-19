package packet

import (
	"errors"
	"fmt"
)

// PacketBytes ensures a coherent naming scheme for packets' internal byte slice
// references.
type PacketBytes struct {
	// Contents contains the entire packet's bytes.
	Contents []byte

	// Payload contains the packet's payload. That is, the contents minus the
	// packet header.
	Payload []byte
}

type Packet interface {
	// Link() returns the parsed link-layer frame.
	Link() Layer

	// Network() returns the parsed network-layer packet.
	Network() Layer

	// Transport() returns the transport-layer segment.
	Transport() Layer
}

// Implementation guard
var _ Packet = new(packet)

// packet contains an instance that satisfies the Packet interface.
type packet struct {
	link    *Ethernet
	network Layer
}

// Decode copies the input bytes, and eagerly decodes the provided byte slice.
func Decode(b []byte) (Packet, error) {
	// Copy input bytes
	cpy := make([]byte, len(b))
	copy(cpy, b)
	b = cpy

	var p packet
	eth := new(Ethernet)
	if err := eth.Unmarshal(b); err != nil {
		return nil, err
	}

	p.link = eth
	if err := p.decodeEthernetFrame(eth); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *packet) decodeEthernetFrame(eth *Ethernet) error {
	switch eth.EthernetType {
	case EthernetTypeARP:
		arp := new(ARP)
		if err := arp.Unmarshal(eth.Payload); err != nil {
			return err
		}
		p.network = arp
	case EthernetTypeIPv4:
		ip := new(IPv4)
		if err := ip.Unmarshal(eth.Payload); err != nil {
			return err
		}
		p.network = ip
	case EthernetTypeIPv6:
		return errors.New("IPv6 not supported")
	default:
		fmt.Printf("unknown network protocol %d\n", eth.EthernetType)
	}
	return nil
}

func (p packet) Link() Layer {
	return p.link
}

func (p packet) Network() Layer {
	return p.network
}

func (p packet) Transport() Layer {
	return p.network
}
