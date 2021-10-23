package tc

import (
	"testing"

	"github.com/mdlayher/netlink"
)

func skipAttribute(t *testing.T, typ uint16, skip []uint16) bool {
	t.Helper()

	for _, s := range skip {
		if s == typ {
			return true
		}
	}
	return false
}

// stripRateTable is a helper function used only in tests.
func stripRateTable(t *testing.T, orig []byte, skip []uint16) ([]byte, error) {
	t.Helper()

	var attrs []netlink.Attribute

	ad, err := netlink.NewAttributeDecoder(orig)
	if err != nil {
		return []byte{}, err
	}
	for ad.Next() {
		if !skipAttribute(t, ad.Type(), skip) {
			attrs = append(attrs, netlink.Attribute{
				Type: ad.Type(),
				Data: ad.Bytes(),
			})
		}
	}

	return netlink.MarshalAttributes(attrs)
}
