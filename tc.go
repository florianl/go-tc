//+build linux

package tc

import (
	"encoding/binary"
	"unsafe"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// Tc represents a RTNETLINK handler
type Tc struct {
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
func Open(config *Config) (*Tc, error) {
	var tc Tc

	con, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{NetNS: config.NetNS})
	if err != nil {
		return nil, err
	}
	tc.con = con

	return &tc, nil
}

// Close the connection to the netfilter route subsystem
func (tc *Tc) Close() error {
	return tc.con.Close()
}

func (tc *Tc) query(req netlink.Message) ([]netlink.Message, error) {
	verify, err := tc.con.Send(req)
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	return tc.con.Receive()
}

// TcObject represents a generic traffic controll object
type TcObject struct {
	Tcmsg
	TcInfo
}

// TcInfo contains attributes of a queueing discipline
type TcInfo struct {
	Kind         string
	EgressBlock  uint32
	IngressBlock uint32
	HwOffload    uint8
	Chain        uint32
	TcStats      *TcStats
	TcXStats     *TcStats
	TcStats2     *TcStats2
	FqCodel      *TcFqCodel
	BPF          *TcBPF
}

// TcFqCodel contains attributes of the fq_codel discipline
type TcFqCodel struct {
	Target        uint32
	Limit         uint32
	Interval      uint32
	ECN           uint32
	Flows         uint32
	Quantum       uint32
	CEThreshold   uint32
	DropBatchSize uint32
	MemoryLimit   uint32
}

// TcBPF contains attributes of the bpf discipline
type TcBPF struct {
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      string
	ID       uint32
}

// BuildHandle is a simple helper function to construct the handle for the Tcmsg struct
func BuildHandle(major, minor uint16) uint32 {
	return uint32(major)<<16 | uint32(minor)
}
