package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCsumUnspec = iota
	tcaCsumParms
	tcaCsumTm
	tcaCsumPad
)

// Csum contains attributes of the csum discipline
type Csum struct {
	Parms *CsumParms
	Tm    *Tcft
}

// marshalCsum returns the binary encoding of Csum
func marshalCsum(info *Csum) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Csum: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCsumParms, Data: data})
	}
	return marshalAttributes(options)
}

// unmarshalCsum parses the csum-encoded data and stores the result in the value pointed to by info.
func unmarshalCsum(data []byte, info *Csum) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaCsumParms:
			parms := &CsumParms{}
			if err := unmarshalStruct(ad.Bytes(), parms); err != nil {
				return err
			}
			info.Parms = parms
		case tcaCsumTm:
			tcft := &Tcft{}
			if err := unmarshalStruct(ad.Bytes(), tcft); err != nil {
				return err
			}
			info.Tm = tcft
		case tcaCsumPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalCsum()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// CsumParms from include/uapi/linux/tc_act/tc_csum.h
type CsumParms struct {
	Index       uint32
	Capab       uint32
	Action      uint32
	RefCnt      uint32
	BindCnt     uint32
	UpdateFlags uint32
}
