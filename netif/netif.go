package netif

import (
	"errors"
	"net"
)

// FindIF finds a valid network interface
func FindIF() (net.Interface, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, err
	}
	for _, iface := range ifs {
		bm := net.FlagBroadcast | net.FlagMulticast | net.FlagUp
		if iface.Flags&bm == bm {
			return iface, nil
		}
	}
	return net.Interface{}, errors.New("iface not found")
}
