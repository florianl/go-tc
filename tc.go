package tc

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/josharian/native"
	"github.com/mdlayher/netlink"
)

// tcMsgHdrLen is the length in bytes of the Msg header that precedes the
// attributes in every qdisc/class/filter netlink message.
const tcMsgHdrLen = 20

// tcConn defines a subset of netlink.Conn.
type tcConn interface {
	Close() error
	JoinGroup(group uint32) error
	LeaveGroup(group uint32) error
	Receive() ([]netlink.Message, error)
	Send(m netlink.Message) (netlink.Message, error)
	SetOption(option netlink.ConnOption, enable bool) error
	SetReadDeadline(t time.Time) error
}

var _ tcConn = &netlink.Conn{}

// Tc represents a RTNETLINK wrapper
//
// A Tc wraps a single netlink socket. That socket cannot be used for a
// synchronous call (e.g. Qdisc().Get(), Filter().Add()) and Monitor() or
// MonitorWithErrorFunc() at the same time. Replies to synchronous calls and
// messages destined for a Monitor loop are read from the same underlying
// connection with no way to tell them apart, so one call can accidentally
// consume the message meant for the other. Open a dedicated Tc for
// Monitor()/MonitorWithErrorFunc() and use a separate one for synchronous
// calls. Once a Monitor loop is active on a Tc, synchronous calls on it
// return ErrMonitorActive.
type Tc struct {
	con tcConn

	// TODO: Switch to *atomic.Bool once upgraded to Go 1.19.
	monitoring *atomic.Value
}

var nativeEndian = native.Endian

// newMonitoringFlag returns an atomic.Value ready to guard a Tc's connection.
func newMonitoringFlag() *atomic.Value {
	v := new(atomic.Value)
	v.Store(false)
	return v
}

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
	tc.monitoring = newMonitoringFlag()

	return &tc, nil
}

// monitoringActive reports whether a Monitor loop is currently receiving on
// this connection.
func (tc *Tc) monitoringActive() bool {
	return tc.monitoring != nil && tc.monitoring.Load().(bool)
}

// tryStartMonitoring marks this connection as being monitored. It reports
// false if a Monitor loop is already active on it.
func (tc *Tc) tryStartMonitoring() bool {
	if tc.monitoring == nil {
		return true
	}
	return tc.monitoring.CompareAndSwap(false, true)
}

// stopMonitoring clears the monitoring flag set by tryStartMonitoring.
func (tc *Tc) stopMonitoring() {
	if tc.monitoring != nil {
		tc.monitoring.Store(false)
	}
}

// SetOption allows to enable or disable netlink socket options.
func (tc *Tc) SetOption(o netlink.ConnOption, enable bool) error {
	return tc.con.SetOption(o, enable)
}

// Close the connection
func (tc *Tc) Close() error {
	return tc.con.Close()
}

func (tc *Tc) query(req netlink.Message) ([]netlink.Message, error) {
	if tc.monitoringActive() {
		return nil, ErrMonitorActive
	}

	verify, err := tc.con.Send(req)
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	return tc.con.Receive()
}

func (tc *Tc) action(action int, flags netlink.HeaderFlags, msg interface{}, opts []tcOption) error {
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
		switch msg.Header.Type {
		case netlink.Error:
			if len(msg.Data) < 4 {
				return fmt.Errorf("received error message from netlink: %w", ErrShortMsg)
			}
			errCode := bytesToInt32(msg.Data[:4])
			// Check if the success message is embedded encoded as error code 0:
			if errCode != 0 {
				return fmt.Errorf("received error from netlink: %#v", msg)
			}
		case netlink.Overrun:
			return fmt.Errorf("lost netlink data: %#v", msg)
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
		if len(msg.Data) < tcMsgHdrLen {
			return results, fmt.Errorf("received message from netlink: %w", ErrShortMsg)
		}
		var result Object
		if err := unmarshalStruct(msg.Data[:tcMsgHdrLen], &result.Msg); err != nil {
			return results, err
		}
		if err := extractTcmsgAttributes(action, msg.Data[tcMsgHdrLen:], &result.Attribute); err != nil {
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
	EgressBlock  *uint32
	IngressBlock *uint32
	HwOffload    *uint8
	Chain        *uint32
	Stats        *Stats
	XStats       *XStats
	Stats2       *GenStats
	Stab         *Stab
	ExtWarnMsg   string

	// Filters
	Basic    *Basic
	BPF      *Bpf
	Cgroup   *Cgroup
	U32      *U32
	Rsvp     *Rsvp
	Route4   *Route4
	Fw       *Fw
	Flow     *Flow
	Flower   *Flower
	Matchall *Matchall
	TcIndex  *TcIndex

	// Classless qdiscs
	Cake    *Cake
	FqCodel *FqCodel
	FqPie   *FqPie
	Codel   *Codel
	Fq      *Fq
	Pie     *Pie
	Hhf     *Hhf
	Tbf     *Tbf
	Sfb     *Sfb
	Sfq     *Sfq
	Red     *Red
	Gred    *Gred
	MqPrio  *MqPrio
	Pfifo   *FifoOpt
	Bfifo   *FifoOpt
	Choke   *Choke
	Netem   *Netem
	Plug    *Plug

	// Classful qdiscs
	Cbs      *Cbs
	Htb      *Htb
	Hfsc     *Hfsc
	HfscQOpt *HfscQOpt
	Dsmark   *Dsmark
	Drr      *Drr
	Cbq      *Cbq
	Atm      *Atm
	Qfq      *Qfq
	Prio     *Prio
	TaPrio   *TaPrio
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
	FqPie   *FqPieXStats
	Fq      *FqQdStats
	Hfsc    *HfscXStats
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
	} else if v.FqPie != nil {
		return marshalStruct(v.FqPie)
	}
	return []byte{}, fmt.Errorf("could not marshal XStat")
}

// HookFunc is a function, which is called for each altered RTNETLINK Object.
// Return something different than 0, to stop receiving messages.
// action will have the value of unix.RTM_[NEW|GET|DEL][QDISC|TCLASS|FILTER].
type HookFunc func(action uint16, m Object) int

// ErrorFunc is a function that receives all errors that happen while reading
// from a Netlinkgroup. To stop receiving messages return something different than 0.
type ErrorFunc func(e error) int

// MonitorWithErrorFunc handles NETLINK_ROUTE messages and calls for each HookFunc.
// Received errors tigger the given ErrorFunc.
//
// This dedicates tc's connection to monitoring until the loop stops: do not use tc, or any
// Qdisc()/Filter()/Class()/Chain()/Actions() derived from it, for synchronous calls while the
// loop is running, and do not call Monitor()/MonitorWithErrorFunc() again on it. Either will
// return ErrMonitorActive. See the Tc doc comment for details, and use a separate Tc, opened
// via Open(), for synchronous calls.
func (tc *Tc) MonitorWithErrorFunc(ctx context.Context, deadline time.Duration,
	fn HookFunc, errfn ErrorFunc,
) error {
	return tc.monitor(ctx, deadline, fn, errfn)
}

// Monitor NETLINK_ROUTE messages
//
// This dedicates tc's connection to monitoring until the loop stops: do not use tc, or any
// Qdisc()/Filter()/Class()/Chain()/Actions() derived from it, for synchronous calls while the
// loop is running, and do not call Monitor()/MonitorWithErrorFunc() again on it. Either will
// return ErrMonitorActive. See the Tc doc comment for details, and use a separate Tc, opened
// via Open(), for synchronous calls.
//
// Deprecated: Use MonitorWithErrorFunc() instead.
func (tc *Tc) Monitor(ctx context.Context, deadline time.Duration, fn HookFunc) error {
	return tc.monitor(ctx, deadline, fn, func(err error) int {
		if opError, ok := err.(*netlink.OpError); ok {
			if opError.Timeout() || opError.Temporary() {
				return 0
			}
		}
		return 1
	})
}

func (tc *Tc) monitor(ctx context.Context, deadline time.Duration,
	fn HookFunc, errfn ErrorFunc,
) error {
	if !tc.tryStartMonitoring() {
		return ErrMonitorActive
	}

	ifinfomsg, err := marshalStruct(unix.IfInfomsg{
		Family: unix.AF_UNSPEC,
	})
	if err != nil {
		tc.stopMonitoring()
		return err
	}

	rtattr, err := marshalAttributes([]tcOption{
		{Interpretation: vtUint32, Type: unix.IFLA_EXT_MASK, Data: uint32(1)},
	})
	if err != nil {
		tc.stopMonitoring()
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
		tc.stopMonitoring()
		return err
	}

	verify, err := tc.con.Send(req)
	if err != nil {
		tc.con.LeaveGroup(unix.RTNLGRP_TC)
		tc.stopMonitoring()
		return err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		tc.con.LeaveGroup(unix.RTNLGRP_TC)
		tc.stopMonitoring()
		return err
	}

	go func() {
		defer tc.stopMonitoring()
		go func() {
			<-ctx.Done()
			stop := time.Now().Add(deadline)
			tc.con.SetReadDeadline(stop)
			tc.con.LeaveGroup(unix.RTNLGRP_TC)
		}()
		for {
			msgs, err := tc.con.Receive()
			if err != nil {
				if ret := errfn(err); ret != 0 {
					return
				}
				if ctx.Err() != nil {
					return
				}
				continue
			}
			for _, msg := range msgs {
				if len(msg.Data) < tcMsgHdrLen {
					continue
				}
				var monitored Object
				if err := unmarshalStruct(msg.Data[:tcMsgHdrLen], &monitored.Msg); err != nil {
					continue
				}
				if err := extractTcmsgAttributes(int(msg.Header.Type), msg.Data[tcMsgHdrLen:],
					&monitored.Attribute); err != nil {
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
