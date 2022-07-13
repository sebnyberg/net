package packet

// Package packet provides a net.PackageConn implementation for an AF_PACKET
// socket. It is heavily inspired by mdlayher/packet and google/gopacket.
//
// The motivation behind the package is to enable integration testing of frame
// parsers against real-world, raw packets coming into the local network
// interface.
//
