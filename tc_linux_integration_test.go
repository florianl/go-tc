//go:build integration && linux && go1.16
// +build integration,linux,go1.16

package tc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Get requests don't need special priviledges
func TestLinuxTcQueueGet(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	queues, qErr := tcSocket.Qdisc().Get()
	if qErr != nil {
		t.Fatalf("Could not get queues from TC socket: %v", qErr)
	}
	for _, queue := range queues {
		fmt.Printf("%#v\n", queue)
	}
}

func TestLinuxTcFilterGet(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get interfaces: %v", err)
	}

	for _, iface := range ifaces {
		info := Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  0,
			Parent:  HandleIngress,
			Info:    0,
		}
		filters, err := tcSocket.Filter().Get(&info)
		if err != nil {
			t.Fatalf("Could not get filters from TC socket for %s: %v", err, iface.Name)
		}
		for _, filter := range filters {
			fmt.Printf("%s\t%#v\n", iface.Name, filter)
		}
	}
}

func TestLinuxTcChainGet(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get interfaces: %v", err)
	}

	for _, iface := range ifaces {
		info := Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  0,
			Parent:  HandleIngress,
			Info:    0,
		}
		chains, err := tcSocket.Chain().Get(&info)
		if err != nil {
			t.Fatalf("Could not get chains from TC socket for %s: %v", err, iface.Name)
		}
		for _, chain := range chains {
			fmt.Printf("%s\t%#v\n", iface.Name, chain)
		}
	}
}

func TestLinuxTcClassGet(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get interfaces: %v", err)
	}

	for _, iface := range ifaces {
		info := Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  0,
			Parent:  HandleIngress,
			Info:    0,
		}
		classes, err := tcSocket.Class().Get(&info)
		if err != nil {
			t.Fatalf("Could not get class from TC socket for %s: %v", err, iface.Name)
		}
		for _, class := range classes {
			fmt.Printf("%s\t%#v\n", iface.Name, class)
		}
	}
}

func TestSocket(t *testing.T) {
	t.Run("empty Config", func(t *testing.T) {
		tc, err := Open(nil)
		if err != nil {
			t.Fatalf("failed to open netlink socket: %v", err)
		}
		if err = tc.Close(); err != nil {
			t.Fatalf("failed to close test socket: %v", err)
		}
	})
	t.Run("with logger", func(t *testing.T) {
		tc, err := Open(&Config{
			Logger: log.Default(),
		})
		if err != nil {
			t.Fatalf("failed to open netlink socket: %v", err)
		}
		if err = tc.Close(); err != nil {
			t.Fatalf("failed to close test socket: %v", err)
		}
	})
}

func TestMonitorWithErrorFunc(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	hook := func(action uint16, m Object) int {
		fmt.Fprintf(os.Stdout, "[%02d]\t%v\n", action, m)
		return 0
	}

	errFunc := func(err error) int {
		if opError, ok := err.(*netlink.OpError); ok {
			if opError.Timeout() || opError.Temporary() {
				return 0
			}
		}
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := tcSocket.MonitorWithErrorFunc(ctx, 10*time.Millisecond, hook, errFunc); err != nil {
		t.Fatal(err)
	}
	cancel()

	<-ctx.Done()
}

func TestMonitor(t *testing.T) {
	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	hook := func(action uint16, m Object) int {
		fmt.Fprintf(os.Stdout, "[%02d]\t%v\n", action, m)
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := tcSocket.Monitor(ctx, 10*time.Millisecond, hook); err != nil {
		t.Fatal(err)
	}
	cancel()

	<-ctx.Done()
}
