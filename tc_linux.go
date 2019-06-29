//+build linux

package tc

import (
	"golang.org/x/sys/unix"

	"github.com/mdlayher/netlink"
)

// Open establishes a RTNETLINK socket for traffic control
func Open(config *Config) (*Tc, error) {
	var tc Tc

	if config == nil {
		config = &Config{}
	}

	con, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{NetNS: config.NetNS})
	if err != nil {
		return nil, err
	}
	tc.con = con

	return &tc, nil
}

// Close the connection
func (tc *Tc) Close() error {
	return tc.con.Close()
}
