package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaBasicUnspec = iota
	tcaBasicClassID
	tcaBasicEmatches
	tcaBasicAct
	tcaBasicPolice
)

// Basic contains attributes of the basic discipline
type Basic struct {
	ClassID *uint32
	Police  *Police
}

// unmarshalBasic parses the Basic-encoded data and stores the result in the value pointed to by info.
func unmarshalBasic(data []byte, info *Basic) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaBasicPolice:
			pol := &Police{}
			if err := unmarshalPolice(ad.Bytes(), pol); err != nil {
				return err
			}
			info.Police = pol
		case tcaBasicClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalBasic()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalBasic returns the binary encoding of Basic
func marshalBasic(info *Basic) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Basic: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBasicClassID, Data: uint32Value(info.ClassID)})
	}
	return marshalAttributes(options)
}
