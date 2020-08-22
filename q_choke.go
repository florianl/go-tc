package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaChokeUnspec = iota
	tcaChokeParms
	tcaChokeStab
	tcaChokeMaxP
)

// Choke contains attributes of the choke discipline
type Choke struct {
	Parms *RedQOpt
	MaxP  *uint32
}

// unmarshalChoke parses the Choke-encoded data and stores the result in the value pointed to by info.
func unmarshalChoke(data []byte, info *Choke) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaChokeParms:
			opt := &RedQOpt{}
			if err := unmarshalStruct(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		case tcaChokeMaxP:
			info.MaxP = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalChoke()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalChoke returns the binary encoding of Choke
func marshalChoke(info *Choke) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Choke: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaChokeParms, Data: data})
	}

	if info.MaxP != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaChokeMaxP, Data: uint32Value(info.MaxP)})

	}

	return marshalAttributes(options)
}
