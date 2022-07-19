package packet

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"
)

// Interface guard
var _ Layer = new(ARP)

type ARPType uint16

const (
	// ARPTypeNetROM     ARPType = 0
	ARPTypeEther ARPType = 1
	// ARPTypeEEther     ARPType = 2
	// ARPTypeAX25       ARPType = 3
	// ARPTypePRONet     ARPType = 4
	// ARPChaos          ARPType = 5
	// ARPIEEE802        ARPType = 6
	// ARPARCNet         ARPType = 7
	// ARPAppletalk      ARPType = 8
	// ARPFrameRelayDLCI ARPType = 15
	// ARPATM            ARPType = 19
	// ARPSTRIP          ARPType = 23
	// ARPIEEE1394       ARPType = 24
	// ARPEUI64          ARPType = 27
	// ARPINFINIBAND     ARPType = 32
)

type ARPOpCode uint16

const (
	ARPOPCodeRequest ARPOpCode = 1
	ARPOPCodeReply   ARPOpCode = 2
	// ARPOPCodeRRequest  ARPOpCode = 3
	// ARPOPCodeRReply    ARPOpCode = 4
	// ARPOPCodeInRequest ARPOpCode = 8
	// ARPOPCodeInReply   ARPOpCode = 9
	// ARPOPCodeNak       ARPOpCode = 10
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
	PacketBytes
}

func (a *ARP) Unmarshal(data []byte) error {
	if len(data) < 8 {
		return errors.New("ARP too short")
	}
	a.HType = ARPType(binary.BigEndian.Uint16(data[0:2]))
	a.PType = EtherType(binary.BigEndian.Uint16(data[2:4]))
	a.HLen = data[4]
	a.PLen = data[5]
	if len(data) < int(8+a.HLen*2+a.PLen*2) {
		return errors.New("invalid ARP packet len")
	}
	a.Oper = ARPOpCode(binary.BigEndian.Uint16(data[6:8]))
	a.SourceHW = net.HardwareAddr(data[8 : 8+a.HLen])
	var ok bool
	a.SourceIP, ok = netip.AddrFromSlice(data[8+a.HLen : 8+a.HLen+a.PLen])
	if !ok {
		return errors.New("invalid source IP addr")
	}
	a.DestHW = net.HardwareAddr(data[8+a.HLen+a.PLen : 8+a.HLen*2+a.PLen])
	a.DestIP, ok = netip.AddrFromSlice(data[8+a.HLen*2+a.PLen : 8+a.HLen*2+a.PLen*2])
	if !ok {
		return errors.New("invalid source IP addr")
	}
	a.Contents = data
	a.Payload = data[8+a.HLen*2+a.PLen*2:]
	return nil
}
func (e ARP) Type() LayerType {
	return LayerTypeEthernet
}

func (e ARP) GetContents() []byte {
	return e.Contents
}

func (e ARP) GetPayload() []byte {
	return e.Payload
}
