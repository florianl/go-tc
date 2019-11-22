package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCodelUnspec = iota
	tcaCodelTarget
	tcaCodelLimit
	tcaCodelInterval
	tcaCodelECN
	tcaCodelCEThreshold
)

// Codel contains attributes of the codel discipline
type Codel struct {
	Target      uint32
	Limit       uint32
	Interval    uint32
	ECN         uint32
	CEThreshold uint32
}

// unmarshalCodel parses the Codel-encoded data and stores the result in the value pointed to by info.
func unmarshalCodel(data []byte, info *Codel) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaCodelTarget:
			info.Target = ad.Uint32()
		case tcaCodelLimit:
			info.Limit = ad.Uint32()
		case tcaCodelInterval:
			info.Interval = ad.Uint32()
		case tcaCodelECN:
			info.ECN = ad.Uint32()
		case tcaCodelCEThreshold:
			info.CEThreshold = ad.Uint32()
		default:
			return fmt.Errorf("unmarshalCodel()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalCodel returns the binary encoding of Red
func marshalCodel(info *Codel) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Codel: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Target != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelTarget, Data: info.Target})
	}
	if info.Limit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelLimit, Data: info.Limit})
	}
	if info.Interval != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelInterval, Data: info.Interval})
	}
	if info.ECN != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelECN, Data: info.ECN})
	}
	if info.CEThreshold != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCodelCEThreshold, Data: info.CEThreshold})
	}
	return marshalAttributes(options)
}
