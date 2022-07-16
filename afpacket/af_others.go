//go:build !linux
// +build !linux

package afpacket

import (
	"fmt"
	"net"
	"runtime"
	"syscall"
	"time"
)

var errUnimplemented = fmt.Errorf("not implemented for %v", runtime.GOOS)

func listen(_ net.Interface) (*Conn, error)               { return nil, errUnimplemented }
func (c *Conn) readFrom(_ []byte) (int, net.Addr, error)  { return 0, nil, errUnimplemented }
func (c *Conn) writeTo(_ []byte, _ net.Addr) (int, error) { return 0, errUnimplemented }
func (c *Conn) close() error                              { return errUnimplemented }
func (c *Conn) localAddr() net.Addr                       { return nil }
func (c *Conn) setDeadline(t time.Time) error             { return errUnimplemented }
func (c *Conn) setReadDeadline(t time.Time) error         { return errUnimplemented }
func (c *Conn) setWriteDeadline(t time.Time) error        { return errUnimplemented }
