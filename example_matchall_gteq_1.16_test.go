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
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// This example demonstrates how to use an eBPF program in a TC filer/matchall action.
func ExampleMatchall() {
	tcIface := "matchallIface"

	// For the purpose of testing a dummy network interface is set up for the example.
	rtnl, err := setupDummyInterface(tcIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup dummy interface: %v\n", err)
		return
	}
	defer rtnl.Close()

	// Get the net.Interface by its name to which the tc/qdisc and tc/filter with
	// the eBPF program will be attached on.
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

	// Open a netlink/tc connection to the Linux kernel. This connection is
	// used to manage the tc/qdisc and tc/filter to which
	// the eBPF program will be attached
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

	// For the purpose of this example handcraft an eBPF program of type BPF_PROG_TYPE_SCHED_ACT
	// that will be attached to the networking interface via the TC filter/matchall action.
	//
	// Check out ebpf.LoadCollection() and ebpf.LoadCollectionSpec() for different ways to load
	// an eBPF program.
	//
	// For eBPF programs of type BPF_PROG_TYPE_SCHED_ACT the returned code defines the action that
	// will be applied to the network packet. Returning 0 translates to TC_ACT_OK and will terminate
	// the packet processing pipeline within netlink/tc and allows the packet to proceed.
	spec := ebpf.ProgramSpec{
		Name: "matchAll",
		Type: ebpf.SchedACT,
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

	// Create a qdisc/ingress object.
	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(tc.HandleRoot, 0x0),
			Parent:  tc.HandleIngress,
		},
		Attribute: tc.Attribute{
			Kind: "ingress",
		},
	}
	// Attach the qdisc/ingress to the networking interface.
	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign clsact to %s: %v\n", tcIface, err)
		return
	}
	// When deleting the qdisc, the applied filter will also be gone
	defer func() {
		if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete qdisc: %v\n", err)
		}
	}()

	fd := uint32(prog.FD())
	eBpfActionIndex := uint32(73) // Arbitrary ID to reference the bpf action.

	// Create a filter/matchall object.
	filter := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Info:    core.FilterInfo(0, unix.ETH_P_ALL),
			Parent:  tc.HandleIngress + 1,
		},
		Attribute: tc.Attribute{
			Kind: "matchall",
			Matchall: &tc.Matchall{
				Actions: &[]*tc.Action{
					{
						Kind: "bpf",
						Bpf: &tc.ActBpf{
							Parms: &tc.ActBpfParms{
								Index: eBpfActionIndex,
							},
							FD: &fd,
						},
					},
				},
			},
		},
	}

	// Load and attach the filter object to the ingress path of the interface.
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "failed to attach matchall filter: %v\n", err)
		return
	}
}
