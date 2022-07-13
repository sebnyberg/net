//go:build linux
// +build linux

package packet

import (
	"net"
	"time"

	"golang.org/x/sys/unix"
)

func listen(iface net.Interface, socketType int) (*Conn, error) {
	sockfd, err := unix.Socket(
		unix.AF_PACKET,
		unix.SOCK_RAW,
		int(htons(unix.ETH_P_ALL)),
	)
	if err != nil {
		return nil, err
	}
	addr := &unix.SockaddrLinklayer{
		Protocol: htons(unix.ETH_P_ALL),
		Ifindex:  iface.Index,
	}
	err = unix.Bind(sockfd, addr)
	if err != nil {
		return nil, err
	}
	c := Conn{
		sockfd: sockfd,
	}
	return &c, nil
}

func (c *Conn) readFrom(b []byte) (int, net.Addr, error) {
	n, sa, err := unix.Recvfrom(c.sockfd, b, 0)
	if err != nil {
		return n, nil, err
	}
	sall := sa.(*unix.SockaddrLinklayer)
	addr := &Addr{
		network: "packet",
		hwAddr:  net.HardwareAddr(sall.Addr[:sall.Halen]),
	}
	return 0, addr, errUnimplemented
}

func (c *Conn) writeTo(_ []byte, _ net.Addr) (int, error) { return 0, errUnimplemented }
func (c *Conn) close() error                              { return errUnimplemented }
func (c *Conn) localAddr() net.Addr                       { return nil }
func (c *Conn) setDeadline(t time.Time) error             { return errUnimplemented }
func (c *Conn) setReadDeadline(t time.Time) error         { return errUnimplemented }
func (c *Conn) setWriteDeadline(t time.Time) error        { return errUnimplemented }
