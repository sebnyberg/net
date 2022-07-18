package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
)

type IPv4 struct {
	IHL            uint8
	DSCP           uint8
	ECN            uint8
	TotalLen       uint16
	ID             uint16
	Flags          uint8
	FragOffset     uint16
	Hops           uint8
	Proto          uint8
	HeaderChecksum uint16
	Source         netip.Addr
	Destination    netip.Addr
	// Todo: options
	Payload []byte
}

func (p *IPv4) Unmarshal(data []byte) error {
	if len(data) < 20 {
		return errors.New("ip packet too small")
	}
	ver := data[0] >> 4
	if ver != 4 {
		return fmt.Errorf("ip packets must be v4, was %v", ver)
	}
	p.IHL = uint8(data[0] & 0x0F)
	if p.IHL != 5 {
		return errors.New("ip options not supported")
	}
	p.DSCP = uint8(data[1] >> 2)
	p.ECN = uint8(data[1] & 0x03)
	p.TotalLen = binary.BigEndian.Uint16(data[2:4])
	p.ID = binary.BigEndian.Uint16(data[4:6])
	flagsFragOff := binary.BigEndian.Uint16(data[6:8])
	p.Flags = uint8(flagsFragOff >> 13)
	p.FragOffset = flagsFragOff & 0x1FFF
	p.Hops = data[8]
	p.Proto = data[9]
	p.HeaderChecksum = binary.BigEndian.Uint16(data[10:12])
	var ok bool
	p.Source, ok = netip.AddrFromSlice(data[12:16])
	if !ok {
		return errors.New("invalid source ip")
	}
	p.Destination, ok = netip.AddrFromSlice(data[16:20])
	if !ok {
		return errors.New("invalid destination ip")
	}
	p.Payload = data[24:]
	return nil
}
