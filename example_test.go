//+build linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/jsimonetti/rtnetlink"
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
	tcIface := "ExampleQdisc"

	rtnl, err := setupDummyInterface(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup dummy interface: %v\n", err)
		return
	}
	defer rtnl.Close()

	devID, err := net.InterfaceByName(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}
	defer func(devID uint32, rtnl *rtnetlink.Conn) {
		if err := rtnl.Link.Delete(devID); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete interface: %v\n", err)
		}
	}(uint32(devID.Index), rtnl)

	// open a rtnetlink socket
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(0xFFFF, 0x0000),
			Parent:  0xFFFFFFF1,
			Info:    0,
		},
		tc.Attribute{
			Kind: "clsact",
		},
	}

	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to lo: %v\n", err)
		return
	}
	if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not delete clsact qdisc from lo: %v\n", err)
		return
	}
}

// This example demonstrate the use with classic BPF
func Example_cBPF() {
	tcIface := "ExampleClassicBPF"

	rtnl, err := setupDummyInterface(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup dummy interface: %v\n", err)
		return
	}
	defer rtnl.Close()

	devID, err := net.InterfaceByName(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}
	defer func(devID uint32, rtnl *rtnetlink.Conn) {
		if err := rtnl.Link.Delete(devID); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete interface: %v\n", err)
		}
	}(uint32(devID.Index), rtnl)

	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(tc.HandleIngress, 0x0000),
			Parent:  tc.HandleIngress,
			Info:    0,
		},
		tc.Attribute{
			Kind: "clsact",
		},
	}

	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to lo: %v\n", err)
		return
	}
	// when deleting the qdisc, the applied filter will also be gone
	defer tcnl.Qdisc().Delete(&qdisc)

	ops := []byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}
	opsLen := uint16(1)
	classID := uint32(0x1001)
	flags := uint32(0x1)

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  tc.HandleIngress,
			Info:    0x300,
		},
		tc.Attribute{
			Kind: "bpf",
			BPF: &tc.Bpf{
				Ops:     &ops,
				OpsLen:  &opsLen,
				ClassID: &classID,
				Flags:   &flags,
			},
		},
	}
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign cBPF: %v\n", err)
		return
	}
}

func ExampleU32() {
	// example from http://man7.org/linux/man-pages/man8/tc-police.8.html
	tcIface := "ExampleU32"

	rtnl, err := setupDummyInterface(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup dummy interface: %v\n", err)
		return
	}
	defer rtnl.Close()

	devID, err := net.InterfaceByName(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}
	defer func(devID uint32, rtnl *rtnetlink.Conn) {
		if err := rtnl.Link.Delete(devID); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete interface: %v\n", err)
		}
	}(uint32(devID.Index), rtnl)

	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  0xFFFF0000,
			Info:    768,
		},
		tc.Attribute{
			Kind: "u32",
			U32: &tc.U32{
				Sel: &tc.U32Sel{
					Flags: 0x1,
					NKeys: 0x3,
					Keys: []tc.U32Key{
						//  match ip protocol 6 0xff
						{Mask: 0xff0000, Val: 0x60000, Off: 0x800, OffMask: 0x0},
						{Mask: 0xff000f00, Val: 0x5c0, Off: 0x0, OffMask: 0x0},
						{Mask: 0xff0000, Val: 0x100000, Off: 0x2000, OffMask: 0x0},
					},
				},
				Police: &tc.Police{
					Tbf: &tc.Policy{
						Action: 0x1,
						Burst:  0xc35000,
						Rate: tc.RateSpec{
							CellLog:   0x3,
							Linklayer: 0x1,
							CellAlign: 0xffff,
							Rate:      0x1e848,
						},
					},
				},
			},
		},
	}

	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to lo: %v\n", err)
		return
	}
	// when deleting the qdisc, the applied filter will also be gone
	defer tcnl.Qdisc().Delete(&qdisc)

	ops := []byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}
	opsLen := uint16(1)
	classID := uint32(0x1001)
	flags := uint32(0x1)

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  tc.HandleIngress,
			Info:    0x300,
		},
		tc.Attribute{
			Kind: "bpf",
			BPF: &tc.Bpf{
				Ops:     &ops,
				OpsLen:  &opsLen,
				ClassID: &classID,
				Flags:   &flags,
			},
		},
	}
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign cBPF: %v\n", err)
		return
	}
}
