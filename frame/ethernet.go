package frame

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type EthernetType uint16

const (
	EtherTypeTooLow EthernetType = 0x7FFF

	EthernetTypeIPv4 EthernetType = 0x0800
	EthernetTypeARP  EthernetType = 0x0806
	EthernetTypeIPv6 EthernetType = 0x86DD

	EtherTypeTooHigh EthernetType = 0x86DE
)

type Ethernet struct {
	Destination  net.HardwareAddr
	Source       net.HardwareAddr
	EthernetType EthernetType
	Payload      []byte
	// Todo: VLAN
}

func (e *Ethernet) Unmarshal(data []byte) error {
	if len(data) < 14 {
		return errors.New("ethernet packet too small")
	}
	e.Destination = net.HardwareAddr(data[0:6])
	e.Source = net.HardwareAddr(data[6:12])
	e.EthernetType = EthernetType(binary.BigEndian.Uint16(data[12:14]))
	e.Payload = data[14:]
	if e.EthernetType <= EtherTypeTooLow || e.EthernetType >= EtherTypeTooHigh {
		return fmt.Errorf("unknown ether type, %x", e)
	}
	return nil
}
