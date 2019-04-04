//+build linux

package rtnetlink

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

type RtNlQdisc struct {
	RtNl
}

type QdiscHandle struct {
	Major uint16
	Minor uint16
}

const (
	rtm_newqdisc = 36
	rtm_delqdisc = 37
	rtm_getqdisc = 38
)

type Qdisc struct {
	Tcmsg
	QdiscInfo
}

type QdiscInfo struct {
	Kind         string
	EgressBlock  uint32
	IngressBlock uint32
	HwOffload    uint8
	Chain        uint32
	TcStats      *TcStats
	TcXStats     *TcStats
	TcStats2     *TcStats2
	FqCodel      *QdiscFqCodel
	BPF          *QdiscBPF
}

type QdiscFqCodel struct {
	Target        uint32
	Limit         uint32
	Interval      uint32
	ECN           uint32
	Flows         uint32
	Quantum       uint32
	CEThreshold   uint32
	DropBatchSize uint32
	MemoryLimit   uint32
}

type QdiscBPF struct {
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      string
	ID       uint32
}

func (rtnl *RtNl) Qdisc() *RtNlQdisc {
	return &RtNlQdisc{*rtnl}
}

func (qd *RtNlQdisc) action(action int, dev string, handle QdiscHandle, parent uint32, qdiscName string) error {
	devID, err := net.InterfaceByName(dev)
	if err != nil {
		fmt.Println(err)
		return err
	}

	tcminfo, err := tcmsgEncode(Tcmsg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(devID.Index),
		Handle:  (uint32(handle.Major) << 16) | uint32(handle.Minor),
		Parent:  0xFFFFFFF1,
		Info:    0,
	})
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, tcminfo...)

	attrs, err := nestAttributes([]RtNlOption{
		RtNlOption{Interpretation: vtString, Type: TCA_KIND, Data: qdiscName},
	})
	if err != nil {
		return err
	}
	data = append(data, attrs...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Acknowledge | netlink.Excl | netlink.Create,
		},
		Data: data,
	}

	msgs, err := qd.query(req)
	if err != nil {
		return err
	}
	fmt.Println(msgs)

	return ErrNotImplemented
}

func (qd *RtNlQdisc) New(dev string, handle QdiscHandle, parent uint32, qdiscName string) error {
	return qd.action(rtm_newqdisc, dev, handle, parent, qdiscName)
}

func (qd *RtNlQdisc) Del(dev string, handle QdiscHandle, parent uint32, qdiscName string) error {
	return qd.action(rtm_delqdisc, dev, handle, parent, qdiscName)
}

func (qd *RtNlQdisc) Get() ([]Qdisc, error) {
	var results []Qdisc

	tcminfo, err := tcmsgEncode(Tcmsg{})
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(rtm_getqdisc),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := qd.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		result := Qdisc{}
		tcmsgDecode(msg.Data[:20], &result.Tcmsg)
		extractTCMSGAttributes(msg.Data[20:], &result.QdiscInfo)
		results = append(results, result)
	}

	return results, nil
}
