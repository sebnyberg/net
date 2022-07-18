package packet

import "fmt"

type Packet interface {
	// Link() returns the parsed link-layer frame
	Link() any

	// Network() returns the parsed network-layer frame
	Network() any
}

type packet struct {
	link    any
	network any
	payload []byte
}

func Decode(b []byte) (*packet, error) {
	var p packet
	p.payload = b
	eth := new(Ethernet)
	if err := eth.Unmarshal(b); err != nil {
		return nil, err
	}
	p.link = eth
	switch eth.EthernetType {
	case EthernetTypeARP:
		arp := new(ARP)
		if err := arp.Unmarshal(eth.Payload); err != nil {
			return nil, err
		}
		p.network = arp
	case EthernetTypeIPv4:
		ip := new(IPv4)
		if err := ip.Unmarshal(eth.Payload); err != nil {
			return nil, err
		}
		p.network = ip
	default:
		fmt.Println("unknown network protocol", eth.EthernetType)
		return nil, nil
	}
	return &p, nil
}

func (p *packet) Link() any {
	return p.link
}

func (p *packet) Network() any {
	return p.network
}
