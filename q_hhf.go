package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaHhfUnspec = iota
	tcaHhfBacklogLimit
	tcaHhfQuantum
	tcaHhfHHFlowsLimit
	tcaHhfResetTimeout
	tcaHhfAdmitBytes
	tcaHhfEVICTTimeout
	tcaHhfNonHHWeight
)

// Hhf contains attributes of the hhf discipline
type Hhf struct {
	BacklogLimit uint32
	Quantum      uint32
	HHFlowsLimit uint32
	ResetTimeout uint32
	AdmitBytes   uint32
	EVICTTimeout uint32
	NonHHWeight  uint32
}

// unmarshalHhf parses the Hhf-encoded data and stores the result in the value pointed to by info.
func unmarshalHhf(data []byte, info *Hhf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHhfBacklogLimit:
			info.BacklogLimit = ad.Uint32()
		case tcaHhfQuantum:
			info.Quantum = ad.Uint32()
		case tcaHhfHHFlowsLimit:
			info.HHFlowsLimit = ad.Uint32()
		case tcaHhfResetTimeout:
			info.ResetTimeout = ad.Uint32()
		case tcaHhfAdmitBytes:
			info.AdmitBytes = ad.Uint32()
		case tcaHhfEVICTTimeout:
			info.EVICTTimeout = ad.Uint32()
		case tcaHhfNonHHWeight:
			info.NonHHWeight = ad.Uint32()
		default:
			return fmt.Errorf("unmarshalHhf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalHhf returns the binary encoding of Hhf
func marshalHhf(info *Hhf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, nil
	}
	// TODO: improve logic and check combinations
	if info.BacklogLimit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfBacklogLimit, Data: info.BacklogLimit})
	}
	if info.Quantum != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfQuantum, Data: info.Quantum})
	}
	if info.HHFlowsLimit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfHHFlowsLimit, Data: info.HHFlowsLimit})
	}
	if info.ResetTimeout != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfResetTimeout, Data: info.ResetTimeout})
	}
	if info.AdmitBytes != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfAdmitBytes, Data: info.AdmitBytes})
	}
	if info.EVICTTimeout != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfEVICTTimeout, Data: info.EVICTTimeout})
	}
	if info.NonHHWeight != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHhfNonHHWeight, Data: info.NonHHWeight})
	}

	return marshalAttributes(options)
}
