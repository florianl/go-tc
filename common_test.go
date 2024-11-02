package tc

// Collection of functions that are used often in regular tests.

import (
	"testing"

	"github.com/mdlayher/netlink"
)

// injectTcft alters orig and adds a Tcft struct, that can only be returned by
// the kernel, and returns the new data bytes.
func injectTcft(t *testing.T, orig []byte, tcftAttribute uint16) ([]byte, *Tcft) {
	t.Helper()

	tcft := Tcft{
		Install:  12,
		LastUse:  34,
		Expires:  56,
		FirstUse: 78,
	}

	tcftBytes, err := marshalStruct(tcft)
	if err != nil {
		t.Fatalf("Failed to marshal tcft: %v", err)
	}

	newData := injectAttribute(t, orig, tcftBytes, tcftAttribute)

	return newData, &tcft
}

// injectAttribute is a helper function for tests to inject new data for a tcaAttribute
func injectAttribute(t *testing.T, orig, new []byte, tcaAttribute uint16) []byte {
	t.Helper()

	var attrs []netlink.Attribute

	ad, err := netlink.NewAttributeDecoder(orig)
	if err != nil {
		t.Fatalf("Failed to decode attributes: %v", err)
	}
	for ad.Next() {
		attrs = append(attrs, netlink.Attribute{
			Type: ad.Type(),
			Data: ad.Bytes(),
		})
	}

	attrs = append(attrs, netlink.Attribute{
		Type: tcaAttribute,
		Data: new,
	})

	newData, err := netlink.MarshalAttributes(attrs)
	if err != nil {
		t.Fatalf("Failed to marshal attributes: %v", err)
	}

	return newData
}
