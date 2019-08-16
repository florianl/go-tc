package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaActUnspec = iota
	tcaActKind
	tcaActOptions
	tcaActIndex
	tcaActStats
	tcaActPad
	tcaActCookie
)

// Action represents action attributes of various filters and classes
type Action struct {
	Kind   string
	Index  uint32
	Cookie *Cookie
}

// unmarshalAction parses the Action-encoded data and stores the result in the value pointed to by info.
func unmarshalAction(data []byte, info *Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaActKind:
			info.Kind = ad.String()
		case tcaActIndex:
			info.Index = ad.Uint32()
		case tcaActCookie:
			cookie := &Cookie{}
			if err := unmarshalStruct(ad.Bytes(), cookie); err != nil {
				return err
			}
			info.Cookie = cookie
		default:
			return fmt.Errorf("unmarshalAction()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalAction returns the binary encoding of Action
func marshalAction(info *Action) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Action options are missing")
	}

	// TODO: improve logic and check combinations
	if len(info.Kind) > 0 {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaActKind, Data: info.Kind})
	}
	if info.Index != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActIndex, Data: info.Index})
	}
	return marshalAttributes(options)
}

// Cookie is passed from user to the kernel for actions and classifiers
type Cookie struct {
	Data uint8
	Len  uint32
}
