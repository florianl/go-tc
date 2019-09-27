package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaPieUnspec = iota
	tcaPieTarget
	tcaPieLimit
	tcaPieTUpdate
	tcaPieAlpha
	tcaPieBeta
	tcaPieECN
	tcaPieBytemode
)

// Pie contains attributes of the pie discipline
type Pie struct {
	Target   uint32
	Limit    uint32
	TUpdate  uint32
	Alpha    uint32
	Beta     uint32
	ECN      uint32
	Bytemode uint32
}

// unmarshalPie parses the Pie-encoded data and stores the result in the value pointed to by info.
func unmarshalPie(data []byte, info *Pie) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaPieTarget:
			info.Target = ad.Uint32()
		case tcaPieLimit:
			info.Limit = ad.Uint32()
		case tcaPieTUpdate:
			info.TUpdate = ad.Uint32()
		case tcaPieAlpha:
			info.Alpha = ad.Uint32()
		case tcaPieBeta:
			info.Beta = ad.Uint32()
		case tcaPieECN:
			info.ECN = ad.Uint32()
		case tcaPieBytemode:
			info.Bytemode = ad.Uint32()
		default:
			return fmt.Errorf("extractPieOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalPie returns the binary encoding of Qfq
func marshalPie(info *Pie) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, nil
	}

	// TODO: improve logic and check combinations
	if info.Target != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieTarget, Data: info.Target})
	}
	if info.Limit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieLimit, Data: info.Limit})
	}
	if info.TUpdate != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieTUpdate, Data: info.TUpdate})
	}
	if info.Alpha != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieAlpha, Data: info.Alpha})
	}
	if info.Beta != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieBeta, Data: info.Beta})
	}
	if info.ECN != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieECN, Data: info.ECN})
	}
	if info.Bytemode != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPieBytemode, Data: info.Bytemode})
	}
	return marshalAttributes(options)
}
