package netlab_test

import (
	"context"
	"fmt"
	"net/netip"
	"testing"

	"github.com/sebnyberg/net/netlab"
)

func requireNoError(t *testing.T, v error, args ...any) {
	if v == nil {
		return
	}
	var msg string
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			msg = fmt.Sprintf(s, args[1:]...)
		}
	}
	t.Logf("unexpected error%v\n", ", "+msg)
	t.Fail()
}

func TestNetlab(t *testing.T) {
	internet := netlab.Network{
		Name:   "internet",
		Prefix: netip.MustParsePrefix("20.0.0.1/24"),
	}

	n0 := netlab.Node{Name: "n0"}
	n1 := netlab.Node{Name: "n1"}
	eth01 := n0.Attach("eth0/1", &internet)
	eth10 := n1.Attach("eth1/0", &internet)
	requireNoError(t, eth01.Link(eth10))

	eth01.Send(context.Background(), []byte("hi"))
}
