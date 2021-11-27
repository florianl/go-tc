package tc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/netlink"
)

func TestFqCodelXStats(t *testing.T) {
	tests := map[string]struct {
		val  FqCodelXStats
		err1 error
		err2 error
	}{
		"Qdisc":   {val: FqCodelXStats{Type: 0, Qd: &FqCodelQdStats{MaxPacket: 123}}},
		"Class":   {val: FqCodelXStats{Type: 1, Cl: &FqCodelClStats{Deficit: -1}}},
		"Unknown": {val: FqCodelXStats{Type: 2}, err1: ErrInvalidArg},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFqCodelXStats(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := FqCodelXStats{}
			err2 := unmarshalFqCodelXStats(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("FqCodelXStats missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil-marshalFqCodelXStats", func(t *testing.T) {
		_, err := marshalFqCodelXStats(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unknown-unmarshalFqCodelXStats", func(t *testing.T) {
		var buf bytes.Buffer
		unknownType := uint32(2)
		if err := binary.Write(&buf, nativeEndian, unknownType); err != nil {
			t.Fatalf("failed to marshal bytes: %v", err)
		}
		if err := unmarshalFqCodelXStats(buf.Bytes(), &FqCodelXStats{}); !errors.Is(err, ErrInvalidArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestMarshalAndAlignStruct(t *testing.T) {
	// The ConnmarkParam struct holds 22 bytes and therefore is not
	// aligned to rtaAlignTo.
	unaligned := ConnmarkParam{
		Index:   1,
		Capab:   2,
		Action:  3,
		RefCnt:  4,
		BindCnt: 5,
		Zone:    6,
	}
	returned := ConnmarkParam{}

	bytes, err := marshalAndAlignStruct(&unaligned)
	if err != nil {
		t.Fatalf("Failed to marshal and align struct: %v", err)
	}
	if len(bytes)%rtaAlignTo != 0 {
		t.Fatalf("Alignment of struct failed")
	}
	if err := unmarshalStruct(bytes, &returned); err != nil {
		t.Fatalf("Failed to unmarshal struct: %v", err)
	}
	if diff := cmp.Diff(unaligned, returned); diff != "" {
		t.Fatalf("Struct alignment missmatch (-want +got):\n%s", diff)
	}
}

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
