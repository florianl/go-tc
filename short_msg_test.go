package tc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// shortMsgConn is a tcConn that reproduces a malformed or truncated netlink
// message, which used to make the package panic instead of returning an
// error. It answers each request exactly once with the truncated message,
// then errors on any further call so a Monitor loop terminates instead of
// busy-spinning on the same truncated message forever.
type shortMsgConn struct {
	msgType netlink.HeaderType
	data    []byte

	received bool
}

func (c *shortMsgConn) Close() error                             { return nil }
func (c *shortMsgConn) JoinGroup(uint32) error                   { return nil }
func (c *shortMsgConn) LeaveGroup(uint32) error                  { return nil }
func (c *shortMsgConn) SetOption(netlink.ConnOption, bool) error { return nil }
func (c *shortMsgConn) SetReadDeadline(time.Time) error          { return nil }

func (c *shortMsgConn) Send(m netlink.Message) (netlink.Message, error) {
	return m, nil
}

func (c *shortMsgConn) Receive() ([]netlink.Message, error) {
	if c.received {
		return nil, errors.New("no more messages")
	}
	c.received = true
	return []netlink.Message{{
		Header: netlink.Header{Type: c.msgType},
		Data:   c.data,
	}}, nil
}

func TestGetShortMessage(t *testing.T) {
	tcSocket := &Tc{con: &shortMsgConn{
		msgType: unix.RTM_NEWQDISC,
		data:    []byte{0x01, 0x02, 0x03}, // shorter than the 20 byte Msg header
	}}

	if _, err := tcSocket.Qdisc().Get(); !errors.Is(err, ErrShortMsg) {
		t.Fatalf("expected ErrShortMsg, got: %v", err)
	}
}

func TestActionErrorShortMessage(t *testing.T) {
	tcSocket := &Tc{con: &shortMsgConn{
		msgType: netlink.Error,
		data:    []byte{0x01, 0x02}, // shorter than the 4 bytes needed for the error code
	}}

	info := &Object{Msg: Msg{Ifindex: 1}, Attribute: Attribute{Kind: "qfq"}}
	if err := tcSocket.Qdisc().Add(info); !errors.Is(err, ErrShortMsg) {
		t.Fatalf("expected ErrShortMsg, got: %v", err)
	}
}

func TestActionsGetShortMessage(t *testing.T) {
	tcSocket := &Tc{con: &shortMsgConn{
		msgType: unix.RTM_NEWACTION,
		data:    []byte{0x01, 0x02, 0x03}, // shorter than the 4 byte tcaMsg header
	}}

	if _, err := tcSocket.Actions().Get("gact"); !errors.Is(err, ErrShortMsg) {
		t.Fatalf("expected ErrShortMsg, got: %v", err)
	}
}

func TestMonitorShortMessage(t *testing.T) {
	tcSocket := &Tc{con: &shortMsgConn{
		msgType: unix.RTM_NEWQDISC,
		data:    []byte{0x01, 0x02, 0x03}, // shorter than the 20 byte Msg header
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	hookCalled := false
	err := tcSocket.Monitor(ctx, 10*time.Millisecond, func(action uint16, m Object) int {
		hookCalled = true
		return 1
	})
	if err != nil {
		t.Fatalf("could not start tc monitor: %v", err)
	}

	<-ctx.Done()

	if hookCalled {
		t.Fatalf("hook should not have been called for a truncated message")
	}
}
