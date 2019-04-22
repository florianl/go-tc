//+build linux

package tc

import (
	"encoding/binary"
	"unsafe"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// Tc represents a RTNETLINK wrapper
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

// Open establishes a RTNETLINK socket for traffic control
func Open(config *Config) (*Tc, error) {
	var tc Tc

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

func (tc *Tc) action(action int, flags netlink.HeaderFlags, info *Object, opts []tcOption) error {
	tcminfo, err := tcmsgEncode(&info.Msg)
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, tcminfo...)

	attrs, err := marshalAttributes(opts)
	if err != nil {
		return err
	}
	data = append(data, attrs...)
	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Acknowledge | flags,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		_ = msg
	}

	return nil
}

func (tc *Tc) get(action int, i *Msg) ([]Object, error) {
	var results []Object

	tcminfo, err := tcmsgEncode(i)
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		var result Object
		if err := tcmsgDecode(msg.Data[:20], &result.Msg); err != nil {
			return results, nil
		}
		if err := extractTcmsgAttributes(msg.Data[20:], &result.Attribute); err != nil {
			return results, nil
		}
		results = append(results, result)
	}

	return results, nil
}

// Object represents a generic traffic control object
type Object struct {
	Msg
	Attribute
}

// Attribute contains various elements for traffic control
type Attribute struct {
	Kind         string
	EgressBlock  uint32
	IngressBlock uint32
	HwOffload    uint8
	Chain        uint32
	Stats        *Stats
	XStats       *Stats
	Stats2       *Stats2
	FqCodel      *FqCodel
	BPF          *BPF
}

// FqCodel contains attributes of the fq_codel discipline
type FqCodel struct {
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

// BPF contains attributes of the bpf discipline
type BPF struct {
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      []byte
	ID       uint32
	Action   *Action
}

// BPFActionOptions contains various action attributes
type BPFActionOptions struct {
	OpsLen uint16
	Ops    []byte
	Tcft   *Tcft
	FD     uint32
	Name   string
	Act    *ActBpf
}

// ActionStats contains various statistics of a action
type ActionStats struct {
	Basic     *GenStatsBasic
	RateEst   *GenStatsRateEst
	Queue     *GenStatsQueue
	RateEst64 *GenStatsRateEst64
}

// Action describes a Traffic Control action
type Action struct {
	Kind       string
	Statistics *ActionStats
	BPFOptions *BPFActionOptions
}

// BuildHandle is a simple helper function to construct the handle for the Tcmsg struct
func BuildHandle(major, minor uint16) uint32 {
	return uint32(major)<<16 | uint32(minor)
}
