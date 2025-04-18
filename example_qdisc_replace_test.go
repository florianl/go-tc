//go:build linux
// +build linux

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

func ExampleQdisc_Replace() {
	tcIface := "ExampleQdiscReplace"

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

	if err := addFQQdisc(tcnl, uint32(devID.Index)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add fq qdisc: %v\n", err)
		return
	}
	if err := replaceFQQdisc(tcnl, uint32(devID.Index)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to replace fq qdisc: %v\n", err)
		return
	}
}

// tc qdisc add dev ExampleQdiscReplace root handle 1: fq ce_threshold 4ms
func addFQQdisc(tcnl *tc.Tc, ifIndex uint32) error {
	ceThreshold := uint32(4000)
	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: ifIndex,
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
		},
		Attribute: tc.Attribute{
			Kind: "fq",
			Fq: &tc.Fq{
				CEThreshold: &ceThreshold,
			},
		},
	}

	return tcnl.Qdisc().Add(&qdisc)
}

// tc qdisc replace dev ExampleQdiscReplace root handle 1: fq limit 100
func replaceFQQdisc(tcnl *tc.Tc, ifIndex uint32) error {
	limit := uint32(100)
	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: ifIndex,
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
		},
		Attribute: tc.Attribute{
			Kind: "fq",
			Fq: &tc.Fq{
				PLimit: &limit,
			},
		},
	}

	return tcnl.Qdisc().Replace(&qdisc)
}
