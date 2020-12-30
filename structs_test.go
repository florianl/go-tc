package tc

import (
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
		"Qdisc": {val: FqCodelXStats{Type: 0, Qd: &FqCodelQdStats{MaxPacket: 123}}},
		"Class": {val: FqCodelXStats{Type: 1, Cl: &FqCodelClStats{Deficit: -1}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFqCodelXStats(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
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
	t.Run("nil", func(t *testing.T) {
		_, err := marshalFqCodelXStats(nil)
		if !errors.Is(err, ErrNoArg) {
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

	var attrs []netlink.Attribute

	ad, err := netlink.NewAttributeDecoder(orig)
	if err != nil {
		t.Fatalf("Failed to decode attributes: %v", err)
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		attrs = append(attrs, netlink.Attribute{
			Type: ad.Type(),
			Data: ad.Bytes(),
		})
	}

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
	attrs = append(attrs, netlink.Attribute{
		Type: tcftAttribute,
		Data: tcftBytes,
	})

	newData, err := netlink.MarshalAttributes(attrs)
	if err != nil {
		t.Fatalf("Failed to marshal attributes: %v", err)
	}

	return newData, &tcft
}
