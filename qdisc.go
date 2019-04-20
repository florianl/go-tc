package tc

import (
	"fmt"

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
	options, err := validateQdiscObject(rtmNewQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmNewQdisc, info, options)
}

// Del removes a queueing discipline
func (qd *Qdisc) Del(info *Object) error {
	options, err := validateQdiscObject(rtmDelQdisc, info)
	if err != nil {
		return err
	}
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

func validateQdiscObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, fmt.Errorf("Could not set device ID 0")
	}

	if info.Kind != "clsact" {
		return options, ErrNotImplemented
	}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	if info.Stats != nil || info.XStats != nil || info.Stats2 != nil || info.FqCodel != nil || info.BPF != nil {
		return options, ErrNotImplemented
	}

	if info.EgressBlock != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaEgressBlock, Data: info.EgressBlock})
	}
	if info.IngressBlock != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIngressBlock, Data: info.IngressBlock})
	}
	if info.HwOffload != 0 {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: info.HwOffload})
	}
	if info.Chain != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaChain, Data: info.Chain})
	}
	return options, nil
}
