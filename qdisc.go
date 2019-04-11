//+build linux

package tc

import (
	"github.com/mdlayher/netlink"
)

// Qdisc represents the queueing discipline part of traffic controll
type Qdisc struct {
	Tc
}

const (
	rtmNewQdisc = 36
	rtmDelQdisc = 37
	rtmGetQdisc = 38
)

// Qdisc allows to read and alter queueing disciplins from the rtnetlink socket
func (tc *Tc) Qdisc() *Qdisc {
	return &Qdisc{*tc}
}

func (qd *Qdisc) action(action int, info *TcObject) error {
	tcminfo, err := tcmsgEncode(&info.Tcmsg)
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, tcminfo...)

	attrs, err := processAttributes(&info.TcInfo)
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

	for _, msg := range msgs {
		_ = msg
	}

	return nil
}

// New adds a queueing discipline
func (qd *Qdisc) New(info *TcObject) error {
	return qd.action(rtmNewQdisc, info)
}

// Del removess a queueing discipline
func (qd *Qdisc) Del(info *TcObject) error {
	return qd.action(rtmDelQdisc, info)
}

// Get a queueing discipline
func (qd *Qdisc) Get() ([]TcObject, error) {
	var results []TcObject

	tcminfo, err := tcmsgEncode(&Tcmsg{})
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(rtmGetQdisc),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := qd.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		result := TcObject{}
		if err := tcmsgDecode(msg.Data[:20], &result.Tcmsg); err != nil {
			return results, nil
		}
		if err := extractTCMSGAttributes(msg.Data[20:], &result.TcInfo); err != nil {
			return results, nil
		}
		results = append(results, result)
	}

	return results, nil
}

func processAttributes(info *TcInfo) ([]byte, error) {

	options := []rtNlOption{}

	options = append(options, rtNlOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	return nestAttributes(options)
}
