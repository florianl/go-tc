package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaStabUnspec = iota
	tcaStabBase
	tcaStabData
)

// SizeSpec implements tc_sizespec
type SizeSpec struct {
	CellLog   uint8
	SizeLog   uint8
	CellAlign int16
	Overhead  int32
	LinkLayer uint32
	MPU       uint32
	MTU       uint32
	TSize     uint32
}

// Stab contains attributes of a stab
type Stab struct {
	Base *SizeSpec
	Data *[]byte
}

func unmarshalStab(data []byte, stab *Stab) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaStabBase:
			base := &SizeSpec{}
			if err := unmarshalStruct(ad.Bytes(), base); err != nil {
				return err
			}
			stab.Base = base
		case tcaStabData:
			tmp := ad.Bytes()
			stab.Data = &tmp
		default:
			return fmt.Errorf("unmarshalStab()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
