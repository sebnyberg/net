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

// Packet represents a raw packet flowing through the network.
type Packet struct {
	// Link contains the link-layer representation of the packet.
	Link *Ethernet

	// Network contains the network-layer representation of the packet.
	Network Layer
}

// Decode copies the input bytes, and eagerly decodes the provided byte slice.
func Decode(b []byte) (Packet, error) {
	// Copy input bytes
	cpy := make([]byte, len(b))
	copy(cpy, b)
	b = cpy

	var p Packet
	eth := new(Ethernet)
	if err := eth.Unmarshal(b); err != nil {
		return p, err
	}

	p.Link = eth
	if err := p.decodeEthernetFrame(eth); err != nil {
		return p, err
	}

	return p, nil
}

func (p *Packet) decodeEthernetFrame(eth *Ethernet) error {
	switch eth.EthernetType {
	case EthernetTypeARP:
		arp := new(ARP)
		if err := arp.Unmarshal(eth.Payload); err != nil {
			return err
		}
		p.Network = arp
	case EthernetTypeIPv4:
		ip := new(IPv4)
		if err := ip.Unmarshal(eth.Payload); err != nil {
			return err
		}
		p.Network = ip
	case EthernetTypeIPv6:
		return errors.New("IPv6 not supported")
	default:
		fmt.Printf("unknown network protocol %d\n", eth.EthernetType)
	}
	return nil
}
