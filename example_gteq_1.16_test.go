//go:build go1.16 && linux
// +build go1.16,linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
)

// This example demonstrate how to attach an eBPF program with TC to an interface.
func Example_eBPF() {
	tcIface := "ExampleEBPF"

	// For the purpose of testing a dummy network interface is set up for the example.
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
		// As a dummy network interface was set up for this test make sure to
		// remove this link again once this example finished.
		if err := rtnl.Link.Delete(devID); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete interface: %v\n", err)
		}
	}(uint32(devID.Index), rtnl)

	// Open a netlink/tc connection to the Linux kernel.
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

	// For enhanced error messages from the kernel, it is recommended to set
	// option `NETLINK_EXT_ACK`, which is supported since 4.12 kernel.
	//
	// If not supported, `unix.ENOPROTOOPT` is returned.
	if err := tcnl.SetOption(netlink.ExtendedAcknowledge, true); err != nil {
		fmt.Fprintf(os.Stderr, "could not set option ExtendedAcknowledge: %v\n", err)
		return
	}

	// Create a qdisc/clsact object that will be attached to the ingress part
	// of the networking interface.
	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(tc.HandleRoot, 0x0000),
			Parent:  tc.HandleIngress,
			Info:    0,
		},
		Attribute: tc.Attribute{
			Kind: "clsact",
		},
	}

	// Attach the qdisc/clsact to the networking interface.
	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to %s: %v\n", tcIface, err)
		return
	}
	// When deleting the qdisc, the applied filter will also be gone
	defer tcnl.Qdisc().Delete(&qdisc)

	// Handcraft an eBPF program of type BPF_PROG_TYPE_SCHED_CLS that will be attached to
	// the networking interface via qdisc/clsact.
	//
	// For eBPF programs of type BPF_PROG_TYPE_SCHED_CLS the returned code defines the action that
	// will be applied to the network packet. Returning 0 translates to TC_ACT_OK and will terminate
	// the packet processing pipeline within netlink/tc and allows the packet to proceed.
	spec := ebpf.ProgramSpec{
		Name: "test",
		Type: ebpf.SchedCLS,
		Instructions: asm.Instructions{
			// Set exit code to 0
			asm.Mov.Imm(asm.R0, 0),
			asm.Return(),
		},
		License: "GPL",
	}

	// Load the eBPF program into the kernel.
	prog, err := ebpf.NewProgram(&spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load eBPF program: %v\n", err)
		return
	}

	fd := uint32(prog.FD())
	flags := uint32(0x1)

	// Create a tc/filter object that will attach the eBPF program to the qdisc/clsact.
	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  core.BuildHandle(tc.HandleRoot, tc.HandleMinIngress),
			Info:    0x300,
		},
		tc.Attribute{
			Kind: "bpf",
			BPF: &tc.Bpf{
				FD:    &fd,
				Flags: &flags,
			},
		},
	}

	// Attach the tc/filter object with the eBPF program to the qdisc/clsact.
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not attach filter for eBPF program: %v\n", err)
		return
	}
}
