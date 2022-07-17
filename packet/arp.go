package packet

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"
)

type ARPType uint16

const (
	// ARPTypeNetROM     = 0
	ARPTypeEther = 1
	// ARPTypeEEther     = 2
	// ARPTypeAX25       = 3
	// ARPTypePRONet     = 4
	// ARPChaos          = 5
	// ARPIEEE802        = 6
	// ARPARCNet         = 7
	// ARPAppletalk      = 8
	// ARPFrameRelayDLCI = 15
	// ARPATM            = 19
	// ARPSTRIP          = 23
	// ARPIEEE1394       = 24
	// ARPEUI64          = 27
	// ARPINFINIBAND     = 32
)

type ARPOpCode uint16

const (
	ARPOPCodeRequest = 1
	ARPOPCodeReply   = 2
	// ARPOPCodeRRequest  = 3
	// ARPOPCodeRReply    = 4
	// ARPOPCodeInRequest = 8
	// ARPOPCodeInReply   = 9
	// ARPOPCodeNak       = 10
)

type ARP struct {
	HType    ARPType
	PType    EtherType
	HLen     byte
	PLen     byte
	Oper     ARPOpCode
	SourceHW net.HardwareAddr
	SourceIP netip.Addr
	DestHW   net.HardwareAddr
	DestIP   netip.Addr
}

func (a *ARP) Unmarshal(data []byte) error {
	if len(data) != 28 {
		return errors.New("invalid ARP len")
	}
	a.HType = ARPType(binary.BigEndian.Uint16(data[0:2]))
	a.PType = EtherType(binary.BigEndian.Uint16(data[2:4]))
	a.HLen = data[4]
	a.PLen = data[5]
	a.Oper = ARPOpCode(binary.BigEndian.Uint16(data[6:8]))
	a.SourceHW = net.HardwareAddr(data[8:14])
	var ok bool
	a.SourceIP, ok = netip.AddrFromSlice(data[14:18])
	if !ok {
		return errors.New("invalid source IP addr")
	}
	a.DestHW = net.HardwareAddr(data[18:24])
	a.DestIP, ok = netip.AddrFromSlice(data[24:28])
	if !ok {
		return errors.New("invalid source IP addr")
	}
	return nil
}
