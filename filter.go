//+build linux

package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

// Filter represents the filtering part of rtnetlink
type Filter struct {
	Tc
}

const (
	rtmNewFilter = 44
	rtmDelFilter = 45
	rtmGetFilter = 46
)

// Filter allows to read and alter filters
func (tc *Tc) Filter() *Filter {
	return &Filter{*tc}
}

// Add create a new filter
func (f *Filter) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(rtmNewFilter, info)
	if err != nil {
		return err
	}
	return f.action(rtmNewFilter, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a filter. If the node does not exist yet it is created
func (f *Filter) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(rtmNewFilter, info)
	if err != nil {
		return err
	}
	return f.action(rtmNewFilter, netlink.Create, &info.Msg, options)
}

// Delete removes a filter
func (f *Filter) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateFilterObject(rtmDelFilter, info)
	if err != nil {
		return err
	}
	return f.action(rtmDelFilter, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all filters
func (f *Filter) Get(i *Msg) ([]Object, error) {
	if i == nil {
		return []Object{}, ErrNoArg
	}
	return f.get(rtmGetFilter, i)
}

func validateFilterObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, fmt.Errorf("Could not set device ID 0")
	}

	switch info.Kind {
	case "bpf":
		data, err := MarshalBpf(info.BPF)
		if err != nil {
			return options, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})
	case "u32":
		data, err := MarshalU32(info.U32)
		if err != nil {
			return options, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})
	default:
		return options, ErrNotImplemented
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
