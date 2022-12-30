package tc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
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
				if errors.Is(err2, testcase.err2) {
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
