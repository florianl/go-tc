package tc

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
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

func (tc *Tc) action(action int, flags netlink.HeaderFlags, msg *Msg, opts []tcOption) error {
	tcminfo, err := marshalStruct(msg)
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
		if msg.Header.Type == netlink.Error {
			// TODO: validate NLMSG_ERROR - see https://www.infradead.org/~tgr/libnl/doc/core.html#core_errmsg
		}
	}

	return nil
}

func (tc *Tc) get(action int, i *Msg) ([]Object, error) {
	var results []Object

	tcminfo, err := marshalStruct(i)
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
		if err := unmarshalStruct(msg.Data[:20], &result.Msg); err != nil {
			return results, err
		}
		if err := extractTcmsgAttributes(msg.Data[20:], &result.Attribute); err != nil {
			return results, err
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

// Msg represents a Traffic Control Message
type Msg struct {
	Family  uint32
	Ifindex uint32
	Handle  uint32
	Parent  uint32
	Info    uint32
}

// Attribute contains various elements for traffic control
type Attribute struct {
	Kind         string
	EgressBlock  uint32
	IngressBlock uint32
	HwOffload    uint8
	Chain        uint32
	Stats        *Stats
	XStats       *XStats
	Stats2       *Stats2

	// Filters
	Basic  *Basic
	BPF    *Bpf
	U32    *U32
	Rsvp   *Rsvp
	Route4 *Route4
	Fw     *Fw
	Flow   *Flow

	// Classless qdiscs
	FqCodel *FqCodel
	Codel   *Codel
	Fq      *Fq
	Pie     *Pie
	Hhf     *Hhf
	Tbf     *Tbf
	Sfb     *Sfb
	Red     *Red
	MqPrio  *MqPrio
	Pfifo   *FifoOpt
	Bfifo   *FifoOpt
	Choke   *Choke

	// Classful qdiscs
	Htb    *Htb
	Hfsc   *Hfsc
	Dsmark *Dsmark
	Drr    *Drr
	Cbq    *Cbq
	Atm    *Atm
	Qfq    *Qfq
}

// BuildHandle is a simple helper function to construct the handle for the Tcmsg struct
func BuildHandle(maj, min uint32) uint32 {
	return ((maj & handleMajMask) | (min & handleMinMask))
}

// XStats contains further statistics to the TCA_KIND
type XStats struct {
	Sfb     *SfbXStats
	Sfq     *SfqXStats
	Red     *RedXStats
	Choke   *ChokeXStats
	Htb     *HtbXStats
	Cbq     *CbqXStats
	Codel   *CodelXStats
	Hhf     *HhfXStats
	Pie     *PieXStats
	FqCodel *FqCodelXStats
}

func marshalXStats(v XStats) ([]byte, error) {
	if v.Sfb != nil {
		return marshalStruct(v.Sfb)
	} else if v.Sfq != nil {
		return marshalStruct(v.Sfq)
	} else if v.Red != nil {
		return marshalStruct(v.Red)
	} else if v.Choke != nil {
		return marshalStruct(v.Choke)
	} else if v.Htb != nil {
		return marshalStruct(v.Htb)
	} else if v.Cbq != nil {
		return marshalStruct(v.Cbq)
	} else if v.Codel != nil {
		return marshalStruct(v.Codel)
	} else if v.Hhf != nil {
		return marshalStruct(v.Hhf)
	} else if v.Pie != nil {
		return marshalStruct(v.Pie)
	} else if v.FqCodel != nil {
		return marshalFqCodelXStats(v.FqCodel)
	}
	return []byte{}, fmt.Errorf("could not marshal XStat")
}

// HookFunc is a function, which is called for each altered RTNETLINK Object.
// Return something different than 0, to stop receiving messages.
// action will have the value of unix.RTM_[NEW|GET|DEL][QDISC|TCLASS|FILTER].
type HookFunc func(action uint16, m Object) int

// Monitor NETLINK_ROUTE messages
func (tc *Tc) Monitor(ctx context.Context, deadline time.Duration, fn HookFunc) error {
	ifinfomsg, err := marshalStruct(unix.IfInfomsg{
		Family: unix.AF_UNSPEC,
	})
	if err != nil {
		return err
	}

	rtattr, err := marshalAttributes([]tcOption{{Interpretation: vtUint32, Type: unix.IFLA_EXT_MASK, Data: uint32(1)}})
	if err != nil {
		return err
	}

	data := ifinfomsg
	data = append(data, rtattr...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(unix.RTM_GETLINK),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	if err := tc.con.JoinGroup(unix.RTNLGRP_TC); err != nil {
		return err
	}

	verify, err := tc.con.Send(req)
	if err != nil {
		return err
	}
	_ = verify
	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return err
	}

	go func() {
		defer func() {
			tc.con.LeaveGroup(unix.RTNLGRP_TC)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			deadline := time.Now().Add(deadline)
			tc.con.SetReadDeadline(deadline)
			msgs, err := tc.con.Receive()
			if err != nil {
				if oerr, ok := err.(*netlink.OpError); ok {
					if oerr.Timeout() {
						continue
					}
				}
				return
			}
			for _, msg := range msgs {
				var monitored Object
				if err := unmarshalStruct(msg.Data[:20], &monitored.Msg); err != nil {
					// TODO: add to logging
					continue
				}
				if err := extractTcmsgAttributes(msg.Data[20:], &monitored.Attribute); err != nil {
					// TODO: add to logging
					continue
				}
				if fn(uint16(msg.Header.Type), monitored) != 0 {
					return
				}
			}
		}
	}()
	return nil
}
