package tc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mdlayher/netlink"
)

// blockingConn is a tcConn whose Receive blocks until stop() is called. It
// keeps a Monitor loop "active" on purpose, so the monitoring guard around
// synchronous calls can be exercised deterministically instead of racing a
// background goroutine.
type blockingConn struct {
	closed chan struct{}
}

func newBlockingConn() *blockingConn {
	return &blockingConn{closed: make(chan struct{})}
}

func (c *blockingConn) Close() error                             { return nil }
func (c *blockingConn) JoinGroup(uint32) error                   { return nil }
func (c *blockingConn) LeaveGroup(uint32) error                  { return nil }
func (c *blockingConn) SetOption(netlink.ConnOption, bool) error { return nil }
func (c *blockingConn) SetReadDeadline(time.Time) error          { return nil }

func (c *blockingConn) Send(m netlink.Message) (netlink.Message, error) {
	return m, nil
}

func (c *blockingConn) Receive() ([]netlink.Message, error) {
	<-c.closed
	return nil, errors.New("connection closed")
}

func (c *blockingConn) stop() {
	close(c.closed)
}

func TestMonitorBlocksSynchronousCallsAndItself(t *testing.T) {
	conn := newBlockingConn()
	tcSocket := &Tc{con: conn, monitoring: newMonitoringFlag()}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	noop := func(uint16, Object) int { return 0 }
	if err := tcSocket.Monitor(ctx, time.Millisecond, noop); err != nil {
		t.Fatalf("could not start monitor: %v", err)
	}

	// A second Monitor/MonitorWithErrorFunc call on the same connection must
	// fail immediately instead of racing the first loop's Receive().
	if err := tcSocket.Monitor(ctx, time.Millisecond, noop); !errors.Is(err, ErrMonitorActive) {
		t.Fatalf("expected ErrMonitorActive from a second Monitor call, got: %v", err)
	}

	// A synchronous call, even via a wrapper (Qdisc) derived from the same
	// Tc, must fail while the Monitor loop owns the connection.
	if _, err := tcSocket.Qdisc().Get(); !errors.Is(err, ErrMonitorActive) {
		t.Fatalf("expected ErrMonitorActive from Get(), got: %v", err)
	}

	// Unblock Receive() so the Monitor loop exits and releases the connection.
	conn.stop()

	deadline := time.Now().Add(time.Second)
	for tcSocket.monitoringActive() {
		if time.Now().After(deadline) {
			t.Fatal("monitoring flag was never cleared after the Monitor loop stopped")
		}
		time.Sleep(5 * time.Millisecond)
	}

	// Once released, synchronous calls work again.
	if _, err := tcSocket.Qdisc().Get(); errors.Is(err, ErrMonitorActive) {
		t.Fatalf("Get() still blocked after the Monitor loop released the connection: %v", err)
	}
}

func TestOpenInitializesMonitoringFlag(t *testing.T) {
	tcSocket, err := Open(&Config{})
	if err != nil {
		t.Skipf("could not open netlink socket: %v", err)
	}
	defer tcSocket.Close()

	if tcSocket.monitoring == nil {
		t.Fatal("Open() did not initialize the monitoring flag; the ErrMonitorActive guard is disabled")
	}
}
