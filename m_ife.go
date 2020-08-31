package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaIfeUnspec = iota
	tcaIfeParms
	tcaIfeTm
	tcaIfeDMac
	tcaIfeSMac
	tcaIfeType
	tcaIfeMetaList
	tcaIfePad
)

// Ife contains attribute of the ife discipline
type Ife struct {
	Parms *IfeParms
	SMac  *[]byte
	DMac  *[]byte
	Type  *uint16
	Tm    *Tcft
}

// IfeParms from from include/uapi/linux/tc_act/tc_ife.h
type IfeParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	RefCnt  uint32
	BindCnt uint32
	Flags   uint16
}

// marshalIfe returns the binary encoding of Ife
func marshalIfe(info *Ife) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ife: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Parms != nil {
		data, err := marshalAndAlignStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeParms, Data: data})
	}
	if info.SMac != nil {
		// TODO: use constant instead of 6
		if len(*info.SMac) != 6 {
			return []byte{}, fmt.Errorf("invalid length for SMac")
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeSMac, Data: *info.SMac})
	}
	if info.DMac != nil {
		// TODO: use constant instead of 6
		if len(*info.DMac) != 6 {
			return []byte{}, fmt.Errorf("invalid length for DMac")
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIfeDMac, Data: *info.DMac})
	}
	if info.Type != nil {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaIfeType, Data: *info.Type})
	}
	return marshalAttributes(options)
}

// unmarshalIfe parses the ife-encoded data and stores the result in the value pointed to by info.
func unmarshalIfe(data []byte, info *Ife) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaIfeParms:
			parms := &IfeParms{}
			if err := unmarshalStruct(ad.Bytes(), parms); err != nil {
				return err
			}
			info.Parms = parms
		case tcaIfeSMac:
			tmp := ad.Bytes()
			info.SMac = &tmp
		case tcaIfeDMac:
			tmp := ad.Bytes()
			info.DMac = &tmp
		case tcaIfeTm:
			tcft := &Tcft{}
			if err := unmarshalStruct(ad.Bytes(), tcft); err != nil {
				return err
			}
			info.Tm = tcft
		case tcaIfeType:
			tmp := ad.Uint16()
			info.Type = &tmp
		case tcaIfePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalIfe()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
