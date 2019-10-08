package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaDsmarkUnspec = iota
	tcaDsmarkIndices
	tcaDsmarkDefaultIndex
	tcaDsmarkSetTCIndex
	tcaDsmarkMask
	tcaDsmarkValue
)

// Dsmark contains attributes of the dsmark discipline
type Dsmark struct {
	Indices      uint16
	DefaultIndex uint16
	SetTCIndex   bool
	Mask         uint8
	Value        uint8
}

// unmarshalDsmark parses the Dsmark-encoded data and stores the result in the value pointed to by info.
func unmarshalDsmark(data []byte, info *Dsmark) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaDsmarkIndices:
			info.Indices = ad.Uint16()
		case tcaDsmarkDefaultIndex:
			info.DefaultIndex = ad.Uint16()
		case tcaDsmarkSetTCIndex:
			info.SetTCIndex = ad.Flag()
		case tcaDsmarkMask:
			info.Mask = ad.Uint8()
		case tcaDsmarkValue:
			info.Value = ad.Uint8()
		default:
			return fmt.Errorf("UnmarshalDsmark()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalDsmark returns the binary encoding of Qfq
func marshalDsmark(info *Dsmark) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, nil
	}

	// TODO: improve logic and check combinations
	if info.Indices != 0 {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaDsmarkIndices, Data: info.Indices})
	}
	if info.DefaultIndex != 0 {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaDsmarkDefaultIndex, Data: info.DefaultIndex})
	}
	if info.Mask != 0 {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaDsmarkMask, Data: info.Mask})
	}
	if info.Value != 0 {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaDsmarkValue, Data: info.Value})
	}
	if info.SetTCIndex {
		options = append(options, tcOption{Interpretation: vtFlag, Type: tcaDsmarkSetTCIndex})
	}
	return marshalAttributes(options)
}
