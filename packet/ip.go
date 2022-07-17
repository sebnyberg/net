package packet

import (
	"net/netip"
)

type IP struct {
	Destination netip.Addr
	Source      netip.Addr
	Payload     []byte
}

func (e *IP) Unmarshal(data []byte) error {
	return nil
}
