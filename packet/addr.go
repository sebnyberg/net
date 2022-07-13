package packet

import "net"

type Addr struct {
	network string
	hwAddr  net.HardwareAddr
}

func (a *Addr) Network() string {
	return a.network
}

func (a *Addr) String() string {
	return a.hwAddr.String()
}
