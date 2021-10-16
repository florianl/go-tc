package tc

import (
	"bytes"
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nltest"
)

func testConn(t *testing.T) (*Tc, func()) {
	t.Helper()

	var reqCache []netlink.Message

	c := &Tc{
		con: nltest.Dial(func(req []netlink.Message) ([]netlink.Message, error) {
			if len(req) == 0 {
				// skip validation requests
				return []netlink.Message{}, nil
			}
			if diff := cmp.Diff(1, len(req)); diff != "" {
				t.Fatalf("unexpected number of request messages (-want +got):\n%s", diff)
			}

			var altered []byte
			switch req[0].Header.Type {
			case unix.RTM_NEWTCLASS:
				fallthrough
			case unix.RTM_NEWTFILTER:
				fallthrough
			case unix.RTM_NEWQDISC:
				reqCache = req
			case unix.RTM_GETTFILTER:
				fallthrough
			case unix.RTM_GETQDISC:
				altered = alterResponses(t, &reqCache)
			case unix.RTM_DELTFILTER:
				fallthrough
			case unix.RTM_DELQDISC:
				reqCache = []netlink.Message{}
			default:
			}
			emptyMsg := make([]netlink.Message, 0, 1)
			var data []byte
			tcmsg, err := marshalStruct(&Msg{
				Family:  unix.AF_UNSPEC,
				Ifindex: 0,
				Handle:  0xC001,
				Parent:  0xCAFE,
				Info:    0,
			})
			if err != nil {
				t.Fatalf("could not encode dummy Msg{}: %v", err)
			}
			data = append(data, tcmsg...)
			data = append(data, altered...)

			emptyMsg = append(emptyMsg, netlink.Message{
				Header: netlink.Header{
					Sequence: req[0].Header.Sequence,
					PID:      req[0].Header.PID,
				},
				Data: data,
			})

			return emptyMsg, nil
		}),
	}

	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}
	}
}

var _ netlink.Socket = &socket{}

func (c *socket) Close() error                           { return nil }
func (c *socket) SendMessages(m []netlink.Message) error { c.msgs = append(c.msgs, m...); return nil }
func (c *socket) Send(m netlink.Message) error           { c.msgs = append(c.msgs, m); return nil }
func (c *socket) Receive() ([]netlink.Message, error) {
	if len(c.msgs) > 0 {
		var resp []netlink.Message
		for _, msg := range c.msgs {
			if msg.Header.Type == netlink.HeaderType(unix.RTM_GETLINK) && msg.Header.Flags == (netlink.Request|netlink.Dump) {
				resp = append(resp, netlink.Message{
					Header: netlink.Header{
						Type:  unix.RTM_NEWQDISC,
						Flags: netlink.Replace,
					},
					Data: []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf1, 0xff, 0xff, 0xff, 0x01,
						0x00, 0x00, 0x00, 0x0b, 0x00, 0x01, 0x00, 0x63, 0x6c, 0x73, 0x61, 0x63, 0x74, 0x00, 0x00, 0x04, 0x00, 0x02,
						0x00, 0x05, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x07, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x03,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x2c, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				})
			}
		}
		c.msgs = nil
		return resp, nil
	}
	return []netlink.Message{}, nil
}
func (c *socket) JoinGroup(g uint32) error  { return nil }
func (c *socket) LeaveGroup(g uint32) error { return nil }

// A socket is a netlink.Socket used for testing.
type socket struct {
	msgs []netlink.Message
}

func testHookConn(t *testing.T) (*Tc, func()) {
	t.Helper()

	hookSocket := &socket{}
	c := &Tc{con: netlink.NewConn(hookSocket, 1)}

	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}
	}
}

func TestMonitor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tcSocket, done := testHookConn(t)
	defer done()

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	testHook := func(action uint16, m Object) int {
		t.Logf("Action: %d\nObject: %#v\n", action, m)
		return 1
	}

	// the deadline of 10 * time.Millisecond does not have an effect for the test,
	// as the functionality is not implemented for the test socket
	err := tcSocket.Monitor(ctx, 10*time.Millisecond, testHook)
	if err != nil {
		t.Fatalf("could not start tc monitor: %v", err)
	}

	<-ctx.Done()
}

func alterResponses(t *testing.T, cache *[]netlink.Message) []byte {
	t.Helper()
	var tmp []Object
	var dataStream []byte

	// Decode data from cache
	for _, msg := range *cache {
		var result Object
		if err := extractTcmsgAttributes(0xCAFE, msg.Data[20:], &result.Attribute); err != nil {
			t.Fatalf("could not decode attributes: %v", err)
		}
		tmp = append(tmp, result)
	}

	var stats2 bytes.Buffer
	if err := binary.Write(&stats2, nativeEndian, &Stats2{
		Bytes:      42,
		Packets:    1,
		Qlen:       1,
		Backlog:    0,
		Drops:      0,
		Requeues:   0,
		Overlimits: 42,
	}); err != nil {
		t.Fatalf("could not encode stats2: %v", err)
	}

	var stats bytes.Buffer
	if err := binary.Write(&stats, nativeEndian, &Stats{
		Bytes:      32,
		Packets:    1,
		Drops:      0,
		Overlimits: 0,
		Bps:        1,
		Pps:        1,
		Qlen:       1,
		Backlog:    0,
	}); err != nil {
		t.Fatalf("could not encode stats: %v", err)
	}

	// Alter and marshal data
	for _, obj := range tmp {
		var data []byte
		var attrs []tcOption
		attrs = append(attrs, tcOption{Interpretation: vtString, Type: tcaKind, Data: obj.Kind})
		attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaStats2, Data: stats2.Bytes()})
		attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaStats, Data: stats.Bytes()})
		attrs = append(attrs, tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: uint8(0)})

		// add XStats
		switch obj.Kind {
		case "fq_codel":
			data, err := marshalXStats(XStats{FqCodel: &FqCodelXStats{Type: 0, Qd: &FqCodelQdStats{}}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "sfb":
			data, err := marshalXStats(XStats{Sfb: &SfbXStats{EarlyDrop: 1, PenaltyDrop: 2, AvgProb: 42}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "red":
			data, err := marshalXStats(XStats{Red: &RedXStats{Early: 1, PDrop: 2, Other: 3, Marked: 4}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "choke":
			data, err := marshalXStats(XStats{Choke: &ChokeXStats{Early: 1, PDrop: 2, Other: 3, Marked: 4, Matched: 5}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "htb":
			data, err := marshalXStats(XStats{Htb: &HtbXStats{Lends: 1, Borrows: 2, Giants: 3}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "cbq":
			data, err := marshalXStats(XStats{Cbq: &CbqXStats{Borrows: 2}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "codel":
			data, err := marshalXStats(XStats{Codel: &CodelXStats{MaxPacket: 3, LDelay: 5}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "hhf":
			data, err := marshalXStats(XStats{Hhf: &HhfXStats{DropOverlimit: 42}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		case "pie":
			data, err := marshalXStats(XStats{Pie: &PieXStats{Prob: 2, Delay: 4, Dropped: 5}})
			if err != nil {
				t.Fatalf("could not marshal Xstats struct: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: data})
		}

		if isFilter(obj.Kind) {
			data, err := marshalFilterOptions(obj.Kind, &obj)
			if err != nil {
				t.Fatalf("failed to remarshal filter options: %v", err)
			}
			attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})

		}

		marshaled, err := marshalAttributes(attrs)
		if err != nil {
			t.Fatalf("could not marshal attributes: %v", err)
		}
		data = append(data, marshaled...)

		dataStream = append(dataStream, data...)

	}
	return dataStream
}
