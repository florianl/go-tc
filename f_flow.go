package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFlowUnspec = iota
	tcaFlowKeys
	tcaFlowMode
	tcaFlowBaseClass
	tcaFlowRShift
	tcaFlowAddend
	tcaFlowMask
	tcaFlowXOR
	tcaFlowDivisor
	tcaFlowAct
	tcaFlowPolice
	tcaFlowEMatches
	tcaFlowPerTurb
)

// Flow contains attributes of the flow discipline
type Flow struct {
	Keys      uint32
	Mode      uint32
	BaseClass uint32
	RShift    uint32
	Addend    uint32
	Mask      uint32
	XOR       uint32
	Divisor   uint32
	PerTurb   uint32
}

//UnmarshalFlow parses the Flow-encoded data and stores the result in the value pointed to by info.
func UnmarshalFlow(data []byte, info *Flow) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFlowKeys:
			info.Keys = ad.Uint32()
		case tcaFlowMode:
			info.Mode = ad.Uint32()
		case tcaFlowBaseClass:
			info.BaseClass = ad.Uint32()
		case tcaFlowRShift:
			info.RShift = ad.Uint32()
		case tcaFlowAddend:
			info.Addend = ad.Uint32()
		case tcaFlowMask:
			info.Mask = ad.Uint32()
		case tcaFlowXOR:
			info.XOR = ad.Uint32()
		case tcaFlowDivisor:
			info.Divisor = ad.Uint32()
		case tcaFlowPerTurb:
			info.PerTurb = ad.Uint32()
		default:
			return fmt.Errorf("UnmarshalFlow()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// MarshalFlow returns the binary encoding of Bpf
func MarshalFlow(info *Flow) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Flow options are missing")
	}

	// TODO: improve logic and check combinations
	if info.Keys != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowKeys, Data: info.Keys})
	}
	if info.Mode != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowMode, Data: info.Mode})
	}
	if info.BaseClass != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowBaseClass, Data: info.BaseClass})
	}
	if info.RShift != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowRShift, Data: info.RShift})
	}
	if info.Addend != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowAddend, Data: info.Addend})
	}
	if info.Mask != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowMask, Data: info.Mask})
	}
	if info.XOR != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowXOR, Data: info.XOR})
	}
	if info.Divisor != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowDivisor, Data: info.Divisor})
	}
	if info.PerTurb != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFlowPerTurb, Data: info.PerTurb})
	}
	return marshalAttributes(options)
}
