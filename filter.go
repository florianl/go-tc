package tc

import (
	"fmt"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Filter represents the filtering part of rtnetlink
type Filter struct {
	Tc
}

// Filter allows to read and alter filters
func (tc *Tc) Filter() *Filter {
	return &Filter{*tc}
}

// Add create a new filter
func (f *Filter) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(unix.RTM_NEWTFILTER, info)
	if err != nil {
		return err
	}
	return f.action(unix.RTM_NEWTFILTER, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a filter. If the node does not exist yet it is created
func (f *Filter) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(unix.RTM_NEWTFILTER, info)
	if err != nil {
		return err
	}
	return f.action(unix.RTM_NEWTFILTER, netlink.Create, &info.Msg, options)
}

// Delete removes a filter
func (f *Filter) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(unix.RTM_DELTFILTER, info)
	if err != nil {
		return err
	}
	return f.action(unix.RTM_DELTFILTER, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all filters
func (f *Filter) Get(i *Msg) ([]Object, error) {
	if i == nil {
		return []Object{}, ErrNoArg
	}
	return f.get(unix.RTM_GETTFILTER, i)
}

func validateFilterObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, fmt.Errorf("could not set device ID 0")
	}

	var data []byte
	var err error
	switch info.Kind {
	case "basic":
		data, err = marshalBasic(info.Basic)
	case "flow":
		data, err = marshalFlow(info.Flow)
	case "fw":
		data, err = marshalFw(info.Fw)
	case "route4":
		data, err = marshalRoute4(info.Route4)
	case "rsvp":
		data, err = marshalRsvp(info.Rsvp)
	case "bpf":
		data, err = marshalBpf(info.BPF)
	case "u32":
		data, err = marshalU32(info.U32)
	default:
		return options, ErrNotImplemented
	}
	if err != nil {
		return options, err
	}
	if len(data) < 1 {
		if action == unix.RTM_NEWTFILTER {
			return options, ErrNoArg
		}
	} else {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})
	}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	if info.Stats != nil || info.XStats != nil || info.Stats2 != nil || info.FqCodel != nil {
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
