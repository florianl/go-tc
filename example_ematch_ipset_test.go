//go:build linux
// +build linux

package tc_test

import (
	"fmt"
	"net"
	"os"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/internal/unix"
	"github.com/jsimonetti/rtnetlink"
)

func ExampleEmatch_IPSetMatch() {
	tcIface := "ExampleEmatchIpset"

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

	classID := uint32(42)

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(devID.Index),
		},
		tc.Attribute{
			Kind: "basic",
			Basic: &tc.Basic{
				ClassID: &classID,
				Ematch: &tc.Ematch{
					Hdr: &tc.EmatchTreeHdr{
						NMatches: 1,
					},
					Matches: &[]tc.EmatchMatch{
						{
							Hdr: tc.EmatchHdr{
								Kind: tc.EmatchIPSet,
							},
							IPSetMatch: &tc.IPSetMatch{
								IPSetID: 1337,
								Dir:     []tc.IPSetDir{tc.IPSetSrc},
							},
						},
					},
				},
			},
		},
	}

	if err := tcnl.Filter().Add(&filter); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign filter: %v\n", err)
		return
	}
}
