package tc

import (
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
func (f *Tc) Filter() *Filter {
	return &Filter{*f}
}

// New adds a filter
func (f *Filter) New() error {
	return ErrNotImplemented
}

// Del removes a filter
func (f *Filter) Del() error {
	return ErrNotImplemented
}

// Get fetches all filters
func (f *Filter) Get(i *Msg) ([]Object, error) {
	var results []Object

	tcminfo, err := tcmsgEncode(i)
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(rtmGetFilter),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := f.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		var result Object
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
