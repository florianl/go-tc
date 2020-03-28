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

// setupDummyInterface installs a temporary dummy interface
func setupDummyInterface(iface string) (*rtnetlink.Conn, error) {
	con, err := rtnetlink.Dial(nil)
	if err != nil {
		return &rtnetlink.Conn{}, err
	}

	if err := con.Link.New(&rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   unix.ARPHRD_NETROM,
		Index:  0,
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
		Attributes: &rtnetlink.LinkAttributes{
			Name: iface,
			Info: &rtnetlink.LinkInfo{Kind: "dummy"},
		},
	}); err != nil {
		return con, err
	}

	return con, err
}

func addHfscClass(class *tc.Class, devID, maj, min uint32, serviceCurve *tc.ServiceCurve) (*tc.Object, error) {
	hfsc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: devID,
			Handle:  core.BuildHandle(maj, min),
			Parent:  0x10000,
			Info:    0,
		},
		tc.Attribute{
			Kind: "hfsc",
			Hfsc: &tc.Hfsc{
				Rsc: serviceCurve,
				Fsc: serviceCurve,
				Usc: serviceCurve,
			},
		},
	}

	if err := class.Add(&hfsc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign hfsc class: %v\n", err)
		return nil, err
	}
	return &hfsc, nil
}

func ExampleHfsc() {
	var rtnl *rtnetlink.Conn
	var err error

	tcIface := "tcDev"

	if rtnl, err = setupDummyInterface(tcIface); err != nil {
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
			Handle:  0x10000,
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		// tc qdisc add dev tcDev stab linklayer ethernet mtu 1500 root handle 1: hfsc default 3
		// http://man7.org/linux/man-pages/man8/tc-stab.8.html
		tc.Attribute{
			Kind: "hfsc",
			HfscQOpt: &tc.HfscQOpt{
				DefCls: 3,
			},
			Stab: &tc.Stab{
				Base: &tc.SizeSpec{
					CellLog:   0,
					SizeLog:   0,
					CellAlign: 0,
					Overhead:  0,
					LinkLayer: 1,
					MPU:       0,
					MTU:       1500,
					TSize:     0,
				},
			},
		},
	}

	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		fmt.Fprintf(os.Stderr, "could not assign hfsc to %s: %v\n", tcIface, err)
		return
	}
	defer func() {
		if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete htb qdisc of %s: %v\n", tcIface, err)
			return
		}
	}()

	class1, err := addHfscClass(tcnl.Class(), uint32(devID.Index), 0x10000, 0x1, &tc.ServiceCurve{M2: 0x1e848})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add hfsc: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Class().Delete(class1); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete hfsc class of %s: %v\n", tcIface, err)
			return
		}
	}()

	class2, err := addHfscClass(tcnl.Class(), uint32(devID.Index), 0x10000, 0x2, &tc.ServiceCurve{M2: 0x1e848})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add hfsc: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Class().Delete(class2); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete hfsc class of %s: %v\n", tcIface, err)
			return
		}
	}()

	class3, err := addHfscClass(tcnl.Class(), uint32(devID.Index), 0x10000, 0x3, &tc.ServiceCurve{M2: 0x1e848})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add hfsc: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Class().Delete(class3); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete hfsc class of %s: %v\n", tcIface, err)
			return
		}
	}()

	classes, err := tcnl.Class().Get(&tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(devID.Index)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get all classes: %v\n", err)
	}

	for _, class := range classes {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get interface from id %d: %v", qdisc.Ifindex, err)
			return
		}
		fmt.Printf("%20s\t%s\n", iface.Name, class.Kind)
	}
}
