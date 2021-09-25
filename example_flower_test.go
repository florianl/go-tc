//go:build linux
// +build linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
	"github.com/jsimonetti/rtnetlink"
)

func ExampleFlower() {
	tcIface := "ExampleFlower"

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
			Handle:  core.BuildHandle(tc.HandleRoot, 0),
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

	defer func() {
		if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete qdisc from iface (%d): %v\n", devID.Index, err)
			return
		}
	}()

	srcMac, _ := net.ParseMAC("00:00:5e:00:53:01")
	actions := []*tc.Action{
		{
			Kind: "gact",
			Gact: &tc.Gact{
				Parms: &tc.GactParms{
					Action: 2, // action drop
				},
			},
		},
	}

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  0,
			Parent:  tc.HandleIngress + 1,
			Info:    768,
		},
		tc.Attribute{
			Kind: "flower",
			Flower: &tc.Flower{
				KeyEthSrc: &srcMac,
				Actions:   &actions,
			},
		},
	}

	// tc filter add dev ExampleFlower ingress protocol all prio 1 \
	// flower src_mac 00:00:5e:00:53:01 \
	// action gact drop
	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign flower filter to iface (%d): %v\n", devID.Index, err)
		return
	}
}
