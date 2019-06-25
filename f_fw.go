package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFwUnspec = iota
	tcaFwClassID
	tcaFwPolice
	tcaFwInDev
	tcaFwAct
	tcaFwMask
)

// Fw contains attributes of the fw discipline
type Fw struct {
	ClassID uint32
	Police  *Police
	InDev   string
	Mask    uint32
}

//UnmarshalFw parses the Fw-encoded data and stores the result in the value pointed to by info.
func UnmarshalFw(data []byte, info *Fw) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFwClassID:
			info.ClassID = ad.Uint32()
		case tcaFwInDev:
			info.InDev = ad.String()
		case tcaFwMask:
			info.Mask = ad.Uint32()
		case tcaFwPolice:
			pol := &Police{}
			if err := UnmarshalPolice(ad.Bytes(), pol); err != nil {
				return err
			}
			info.Police = pol
		default:
			return fmt.Errorf("UnmarshalFw()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// MarshalFw returns the binary encoding of Fw
func MarshalFw(info *Fw) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Fw options are missing")
	}

	// TODO: improve logic and check combinations
	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFwClassID, Data: info.ClassID})
	}
	if info.Mask != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFwMask, Data: info.Mask})
	}
	if len(info.InDev) > 0 {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaFwInDev, Data: info.InDev})
	}
	if info.Police != nil {
		data, err := MarshalPolice(info.Police)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaFwPolice, Data: data})
	}

	return marshalAttributes(options)
}
