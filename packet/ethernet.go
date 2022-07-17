package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type EtherType uint16

const (
	EtherTypeTooLow EtherType = 0x7FFF

	EthernetTypeIPv4 EtherType = 0x0800
	EthernetTypeARP  EtherType = 0x0806
	EthernetTypeIPv6 EtherType = 0x86DD

	EtherTypeTooHigh EtherType = 0x86DE
)

type Ethernet struct {
	Destination  net.HardwareAddr
	Source       net.HardwareAddr
	EthernetType EtherType
	Payload      []byte
}

func (e *Ethernet) Unmarshal(data []byte) error {
	if len(data) < 14 {
		return errors.New("ethernet packet too small")
	}
	e.Destination = net.HardwareAddr(data[0:6])
	e.Source = net.HardwareAddr(data[6:12])
	e.EthernetType = EtherType(binary.BigEndian.Uint16(data[12:14]))
	e.Payload = data[14:]
	if e.EthernetType <= EtherTypeTooLow || e.EthernetType >= EtherTypeTooHigh {
		return fmt.Errorf("unknown ether type, %x", e)
	}
	return nil
}
