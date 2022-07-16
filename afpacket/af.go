package afpacket

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

type Conn struct {
	sockfd int
}

func Listen(iface net.Interface) (*Conn, error) {
	return listen(iface)
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
