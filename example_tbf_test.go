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

func ExampleTbf() {
	tcIface := "tcExampleTbf"

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

	linklayerEthernet := uint8(1)
	burst := uint32(0x500000)

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(tc.HandleRoot, 0x0),
			Parent:  tc.HandleRoot,
			Info:    0,
		},

		tc.Attribute{
			Kind: "tbf",
			Tbf: &tc.Tbf{
				Parms: &tc.TbfQopt{
					Mtu:   1514,
					Limit: 0x5000,
					Rate: tc.RateSpec{
						Rate:      0x7d00,
						Linklayer: linklayerEthernet,
						CellLog:   0x3,
					},
				},
				Burst: &burst,
			},
		},
	}

	// tc qdisc add dev tcExampleTbf root tbf burst 20480 limit 20480 mtu 1514 rate 32000bps
	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign tbf to %s: %v\n", tcIface, err)
		return
	}
	defer func() {
		if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete tbf qdisc of %s: %v\n", tcIface, err)
			return
		}
	}()

	qdiscs, err := tcnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get all qdiscs: %v\n", err)
		//return
	}

	fmt.Println("## qdiscs:")
	for _, qdisc := range qdiscs {

		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get interface from id %d: %v", qdisc.Ifindex, err)
			return
		}
		fmt.Printf("%20s\t%-11s\n", iface.Name, qdisc.Kind)
	}
}
