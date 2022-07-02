package tc

// Collection of functions that are used often in regular tests.

import (
	"bytes"
	"encoding/binary"
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

// changeEndianess modifies the endianes for given attributes.
// The kernel expects some arguments to be in a specific endian format but
// returns them in host endian.
func changeEndianess(t *testing.T, orig []byte, attrs map[uint16]valueType) []byte {
	t.Helper()

	var newAttrs []netlink.Attribute

	ad, err := netlink.NewAttributeDecoder(orig)
	if err != nil {
		t.Fatalf("Failed to decode attributes: %v", err)
	}
	for ad.Next() {
		vT, ok := attrs[ad.Type()]
		if !ok {
			newAttrs = append(newAttrs, netlink.Attribute{
				Type: ad.Type(),
				Data: ad.Bytes(),
			})
			continue
		}
		switch vT {
		case vtUint16Be:
			data := bytes.NewBuffer(make([]byte, 0, 2))
			if err := binary.Write(data, binary.BigEndian, ad.Uint16()); err != nil {
				t.Fatalf("changeEndianess for %d: %v", ad.Type(), err)
			}
			newAttrs = append(newAttrs, netlink.Attribute{
				Type: ad.Type(),
				Data: data.Bytes(),
			})
		case vtUint32Be:
			data := bytes.NewBuffer(make([]byte, 0, 4))
			if err := binary.Write(data, binary.BigEndian, ad.Uint32()); err != nil {
				t.Fatalf("changeEndianess for %d: %v", ad.Type(), err)
			}
			newAttrs = append(newAttrs, netlink.Attribute{
				Type: ad.Type(),
				Data: data.Bytes(),
			})
		case vtInt16Be:
			data := bytes.NewBuffer(make([]byte, 0, 2))
			if err := binary.Write(data, binary.BigEndian, ad.Int16()); err != nil {
				t.Fatalf("changeEndianess for %d: %v", ad.Type(), err)
			}
			newAttrs = append(newAttrs, netlink.Attribute{
				Type: ad.Type(),
				Data: data.Bytes(),
			})
		default:
			t.Fatalf("Unexpected valueType %d", vT)
		}
	}

	newData, err := netlink.MarshalAttributes(newAttrs)
	if err != nil {
		t.Fatalf("Failed to marshal attributes: %v", err)
	}

	return newData
}
