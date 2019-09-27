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

// Add creates a new queueing discipline
func (qd *Qdisc) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(rtmNewQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmNewQdisc, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a queueing discipline. If the node does not exist yet it is created
func (qd *Qdisc) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(rtmNewQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmNewQdisc, netlink.Create|netlink.Replace, &info.Msg, options)
}

// Link performs a replace on an existing queueing discipline
func (qd *Qdisc) Link(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(rtmNewQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmNewQdisc, netlink.Replace, &info.Msg, options)
}

// Delete removes a queueing discipline
func (qd *Qdisc) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(rtmDelQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmDelQdisc, netlink.HeaderFlags(0), &info.Msg, options)
}

// Change modifies a queueing discipline 'in place'
func (qd *Qdisc) Change(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(rtmNewQdisc, info)
	if err != nil {
		return err
	}
	return qd.action(rtmNewQdisc, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all queueing disciplines
func (qd *Qdisc) Get() ([]Object, error) {
	return qd.get(rtmGetQdisc, &Msg{})
}

func validateQdiscObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, fmt.Errorf("could not set device ID 0")
	}

	// TODO: improve logic and check combinations
	var data []byte
	var err error
	switch info.Kind {
	case "choke":
		data, err = marshalChoke(info.Choke)
	case "pfifo":
		data, err = marshalStruct(info.Pfifo)
	case "bfifo":
		data, err = marshalStruct(info.Bfifo)
	case "tbf":
		data, err = marshalTbf(info.Tbf)
	case "sfb":
		data, err = marshalSfb(info.Sfb)
	case "red":
		data, err = marshalRed(info.Red)
	case "qfq":
		data, err = marshalQfq(info.Qfq)
	case "pie":
		data, err = marshalPie(info.Pie)
	case "mqprio":
		data, err = marshalMqPrio(info.MqPrio)
	case "hhf":
		data, err = marshalHhf(info.Hhf)
	case "hfsc":
		data, err = marshalHfsc(info.Hfsc)
	case "fq":
		data, err = marshalFq(info.Fq)
	case "dsmark":
		data, err = marshalDsmark(info.Dsmark)
	case "drr":
		data, err = marshalDrr(info.Drr)
	case "codel":
		data, err = marshalCodel(info.Codel)
	case "cbq":
		data, err = marshalCbq(info.Cbq)
	case "atm":
		data, err = marshalAtm(info.Atm)
	case "fq_codel":
		data, err = marshalFqCodel(info.FqCodel)
	case "htb":
		data, err = marshalHtb(info.Htb)
	case "clsact":
		// clsact is parameterless
	case "ingress":
		// ingress is parameterless
	default:
		return options, ErrNotImplemented
	}
	if err != nil {
		return options, err
	}
	if len(data) < 1 && action == rtmNewQdisc {
		if info.Kind != "clsact" && info.Kind != "ingress" {
			return options, ErrNoArg
		}
	} else {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})
	}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	if info.Stats != nil || info.XStats != nil || info.Stats2 != nil || info.BPF != nil {
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
