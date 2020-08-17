// +build linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
)

func ExampleFqCodel() {
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

	target := uint32(0xbb8)
	limit := uint32(0x7d0)
	interval := uint32(0x9c40)
	ecn := uint32(0x0)

	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		tc.Attribute{
			Kind: "fq_codel",
			// http://man7.org/linux/man-pages/man8/tc-fq_codel.8.html
			// fq_codel limit 2000 target 3ms interval 40ms noecn
			FqCodel: &tc.FqCodel{
				Target:   &target,
				Limit:    &limit,
				Interval: &interval,
				ECN:      &ecn,
			},
		},
	}

	if err := rtnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign htb to lo: %v\n", err)
		return
	}
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
