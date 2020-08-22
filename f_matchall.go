package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaMatchallUnspec = iota
	tcaMatchallClassID
	tcaMatchallAct
	tcaMatchallFlags
	tcaMatchallPcnt
	tcaMatchallPad
)

// Matchall contains attributes of the matchall discipline
type Matchall struct {
	ClassID *uint32
	Actions *[]*Action
	Flags   *uint32
}

func unmarshalMatchall(data []byte, info *Matchall) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaMatchallClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaMatchallAct:
			actions := &[]*Action{}
			if err := unmarshalActions(ad.Bytes(), actions); err != nil {
				return err
			}
			info.Actions = actions
		case tcaMatchallFlags:
			info.Flags = uint32Ptr(ad.Uint32())
		case tcaMatchallPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalMatchall()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalMatchall returns the binary encoding of Matchall
func marshalMatchall(info *Matchall) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Matchall: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaMatchallClassID, Data: uint32Value(info.ClassID)})
	}
	if info.Actions != nil {
		data, err := marshalActions(*info.Actions)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaMatchallAct, Data: data})
	}

	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaMatchallFlags, Data: uint32Value(info.Flags)})
	}

	return marshalAttributes(options)
}
