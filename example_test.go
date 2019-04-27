package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"golang.org/x/sys/unix"
)

// This example demonstrate how Get() can get used to read information
func ExampleQdisc_Get() {
	// open a rtnetlink socket
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	qdiscs, err := rtnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get qdiscs: %v\n", err)
		return
	}
	for _, qdisc := range qdiscs {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get interface from id %d: %v", qdisc.Ifindex, err)
			return
		}
		fmt.Printf("%20s\t%s\n", iface.Name, qdisc.Kind)
	}
}

// This example demonstraces how to add a qdisc to an interface and delete it again
func ExampleQdisc() {
	// open a rtnetlink socket
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	devID, err := net.InterfaceByName("lo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  tc.BuildHandle(0xFFFF, 0x0000),
			Parent:  0xFFFFFFF1,
			Info:    0,
		},
		tc.Attribute{
			Kind: "clsact",
		},
	}

	if err := rtnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to lo: %v\n", err)
		return
	}
	if err := rtnl.Qdisc().Delete(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not delete clsact qdisc from lo: %v\n", err)
		return
	}
}

// This example demonstrate the use with classic BPF
func Example_cBPF() {
	// open a rtnetlink socket
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	devID, err := net.InterfaceByName("lo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  tc.BuildHandle(0xFFFF, 0x0000),
			Parent:  0xFFFFFFF1,
			Info:    0,
		},
		tc.Attribute{
			Kind: "clsact",
		},
	}

	if err := rtnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to lo: %v\n", err)
		return
	}
	// when deleting the qdisc, the applied filter will also be gone
	defer rtnl.Qdisc().Delete(&qdisc)

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  tc.Ingress,
			Info:    0x300,
		},
		tc.Attribute{
			Kind: "bpf",
			BPF: &tc.BPF{
				Ops:     []byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff},
				OpsLen:  0x1,
				ClassID: 0x10001,
				Flags:   0x1,
			},
		},
	}
	if err := rtnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign cBPF: %v\n", err)
		return
	}
}
