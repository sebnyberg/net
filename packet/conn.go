package packet

// Package raw provides a net.PackageConn implementation for an AF_PACKET
// socket. It is heavily inspired by mdlayher/packet and google/gopacket.
// My thought behind this packet is to somehow bind a virtual ethernet device to
// the gateway of my virtual in-memory network. That way the network could act
// like an ingress of sorts.
//

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

// errUnimplemented is returned by all functions on non-Linux platforms.
var errUnimplemented = fmt.Errorf("packet: not implemented on %s", runtime.GOOS)

// htons converts a short (uint16) from host-to-network byte order.
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

type Conn struct {
	sockfd int
}

func Listen(iface net.Interface, socketType int) (*Conn, error) {
	return listen(iface, socketType)
}

func (c *Conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.readFrom(p)
}

func (c *Conn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return 0, errUnimplemented
}

func (c *Conn) Close() error {
	return errUnimplemented
}

func (c *Conn) LocalAddr() net.Addr {
	return nil
}

func (c *Conn) SetDeadline(t time.Time) error {
	return errUnimplemented
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return errUnimplemented
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return errUnimplemented
}
