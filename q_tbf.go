package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaTbfUnspec = iota
	tcaTbfParms
	tcaTbfRtab
	tcaTbfPtab
	tcaTbfRate64
	tcaTbfPrate64
	tcaTbfBurst
	tcaTbfPburst
	tcaTbfPad
)

// Tbf contains attributes of the TBF discipline
type Tbf struct {
	Parms   *TbfQopt
	Rtab    []byte
	Ptab    []byte
	Rate64  uint64
	Prate64 uint64
	Burst   uint32
	Pburst  uint32
}

// unmarshalTbf parses the FqCodel-encoded data and stores the result in the value pointed to by info.
func unmarshalTbf(data []byte, info *Tbf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaTbfParms:
			qopt := &TbfQopt{}
			if err := unmarshalStruct(ad.Bytes(), qopt); err != nil {
				return err
			}
			info.Parms = qopt
		case tcaTbfRtab:
			info.Rtab = ad.Bytes()
		case tcaTbfPtab:
			info.Ptab = ad.Bytes()
		case tcaTbfRate64:
			info.Rate64 = ad.Uint64()
		case tcaTbfPrate64:
			info.Prate64 = ad.Uint64()
		case tcaTbfBurst:
			info.Burst = ad.Uint32()
		case tcaTbfPburst:
			info.Pburst = ad.Uint32()
		case tcaTbfPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalTbf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalTbf returns the binary encoding of Tbf
func marshalTbf(info *Tbf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, nil
	}

	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTbfParms, Data: data})
	}
	if info.Rate64 != 0 {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaTbfRate64, Data: info.Rate64})
	}
	if info.Prate64 != 0 {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaTbfPrate64, Data: info.Prate64})
	}
	if info.Burst != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTbfBurst, Data: info.Burst})
	}
	if info.Pburst != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTbfPburst, Data: info.Pburst})
	}
	return marshalAttributes(options)
}

// TbfQopt from include/uapi/linux/pkt_sched.h
type TbfQopt struct {
	Rate     RateSpec
	PeakRate RateSpec
	Limit    uint32
	Buffer   uint32
	Mtu      uint32
}
