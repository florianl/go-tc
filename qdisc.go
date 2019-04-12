package tc

import (
	"github.com/mdlayher/netlink"
)

// Qdisc represents the queueing discipline part of traffic control
type Qdisc struct {
	Tc
}

const (
	rtmNewQdisc = 36
	rtmDelQdisc = 37
	rtmGetQdisc = 38
)

// Qdisc allows to read and alter queues
func (tc *Tc) Qdisc() *Qdisc {
	return &Qdisc{*tc}
}

// New adds a queueing discipline
func (qd *Qdisc) New(info *Object) error {
	options := []tcOption{}

	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	return qd.action(rtmNewQdisc, info, options)
}

// Del removes a queueing discipline
func (qd *Qdisc) Del(info *Object) error {
	options := []tcOption{}

	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	return qd.action(rtmDelQdisc, info, options)
}

// Get fetches all queueing disciplines
func (qd *Qdisc) Get() ([]Object, error) {
	var results []Object

	tcminfo, err := tcmsgEncode(&Msg{})
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
		result := Object{}
		if err := tcmsgDecode(msg.Data[:20], &result.Msg); err != nil {
			return results, nil
		}
		if err := extractTcmsgAttributes(msg.Data[20:], &result.Attribute); err != nil {
			return results, nil
		}
		results = append(results, result)
	}

	return results, nil
}
