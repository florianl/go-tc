//+build linux

package rtnetlink

import (
	"encoding/binary"
	"unsafe"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// RtNl represents a RTNETLINK handler
type RtNl struct {
	con *netlink.Conn
}

// for detailes see https://github.com/tensorflow/tensorflow/blob/master/tensorflow/go/tensor.go#L488-L505
var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
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
