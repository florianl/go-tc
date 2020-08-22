package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMqPrioUnspec = iota
	tcaMqPrioMode
	tcaMqPrioShaper
	tcaMqPrioMinRate64
	tcaMqPrioMaxRate64
)

// MqPrio contains attributes of the mqprio discipline
type MqPrio struct {
	Mode      *uint16
	Shaper    *uint16
	MinRate64 *uint64
	MaxRate64 *uint64
}

// unmarshalMqPrio parses the MqPrio-encoded data and stores the result in the value pointed to by info.
func unmarshalMqPrio(data []byte, info *MqPrio) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaMqPrioMode:
			info.Mode = uint16Ptr(ad.Uint16())
		case tcaMqPrioShaper:
			info.Shaper = uint16Ptr(ad.Uint16())
		case tcaMqPrioMinRate64:
			info.MinRate64 = uint64Ptr(ad.Uint64())
		case tcaMqPrioMaxRate64:
			info.MaxRate64 = uint64Ptr(ad.Uint64())
		default:
			return fmt.Errorf("unmarshalMqPrio()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalMqPrio returns the binary encoding of MqPrio
func marshalMqPrio(info *MqPrio) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("MqPrio: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Mode != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaMqPrioMode, Data: uint16Value(info.Mode)})
	}
	if info.Shaper != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaMqPrioShaper, Data: uint16Value(info.Shaper)})
	}
	if info.MinRate64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaMqPrioMinRate64, Data: uint64Value(info.MinRate64)})
	}
	if info.MaxRate64 != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaMqPrioMaxRate64, Data: uint64Value(info.MaxRate64)})
	}
	return marshalAttributes(options)
}
