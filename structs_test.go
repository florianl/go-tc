package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/netlink"
)

func TestExtractFqCodelXStats(t *testing.T) {
	tests := map[string]struct {
		data     []byte
		expected *FqCodelXStats
		err      error
	}{
		"Qdisc": {data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: &FqCodelXStats{Type: 0, Qd: &FqCodelQdStats{}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			stats := &FqCodelXStats{}
			if err := extractFqCodelXStats(testcase.data, stats); err != nil {
				if testcase.err != nil && testcase.err.Error() == err.Error() {
					// we received the expected error. everything is fine
					return
				}
				t.Fatalf("received error '%v', but expected '%v'", err, testcase.err)
			}
			if diff := cmp.Diff(stats, testcase.expected); diff != "" {
				t.Fatalf("TestExtractFqCodelXStats missmatch (-want +got):\n%s", diff)
			}
		})
	}
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
