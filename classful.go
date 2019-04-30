package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaQfqUnspec = iota
	tcaQfqWeight
	tcaQfqLmax
)

// Qfq contains attributes of the qfq discipline
type Qfq struct {
	Weight uint32
	Lmax   uint32
}

func extractQfqOptions(data []byte, info *Qfq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaQfqWeight:
			info.Weight = ad.Uint32()
		case tcaQfqLmax:
			info.Lmax = ad.Uint32()
		default:
			return fmt.Errorf("extractQfqOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
