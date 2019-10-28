// +build linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"golang.org/x/sys/unix"
)

func ExampleHtb() {
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
			Handle:  tc.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		// configure a very basic hierarchy token bucket (htb) qdisc
		tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Init: &tc.HtbGlob{
					Version:      0x3,
					Rate2Quantum: 0xa,
				},
			},
		},
	}

	if err := rtnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign htb to lo: %v\n", err)
		return
	}
	// delete the qdisc, if this program terminates
	defer func() {
		if err := rtnl.Qdisc().Delete(&qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete htb qdisc of lo: %v\n", err)
			return
		}
	}()
	qdiscs, err := rtnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get all qdiscs: %v\n", err)
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
