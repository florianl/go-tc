//+build linux

package rtnetlink

import (
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// RtNl represents a RTNETLINK handler
type RtNl struct {
	con *netlink.Conn
}

// Open establishes a socket RTNETLINK socket
func Open(config *Config) (*RtNl, error) {
	var rtnl RtNl

	con, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{NetNS: config.NetNS})
	if err != nil {
		return nil, err
	}
	rtnl.con = con

	return &rtnl, nil
}

// Close the connection to the netfilter route subsystem
func (rtnl *RtNl) Close() error {
	return rtnl.con.Close()
}

func (rtnl *RtNl) query(req netlink.Message) ([]netlink.Message, error) {
	verify, err := rtnl.con.Send(req)
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	return rtnl.con.Receive()
}
