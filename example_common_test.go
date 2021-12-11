//go:build linux
// +build linux

package tc_test

// Collection of functions that are used often in examples.

import (
	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// setupDummyInterface installs a temporary dummy interface
func setupDummyInterface(iface string) (*rtnetlink.Conn, error) {
	con, err := rtnetlink.Dial(nil)
	if err != nil {
		return &rtnetlink.Conn{}, err
	}

	if err := con.Link.New(&rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   unix.ARPHRD_NETROM,
		Index:  0,
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
		Attributes: &rtnetlink.LinkAttributes{
			Name: iface,
			Info: &rtnetlink.LinkInfo{Kind: "dummy"},
		},
	}); err != nil {
		return con, err
	}

	return con, err
}
