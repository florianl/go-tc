package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasic(t *testing.T) {
	tests := map[string]struct {
		val  Basic
		err1 error
		err2 error
	}{
		"simple": {val: Basic{ClassID: uint32Ptr(2)}},
		"with ematch & police": {val: Basic{ClassID: uint32Ptr(3),
			Ematch: &Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1},
				Matches: &[]EmatchMatch{
					{Hdr: EmatchHdr{MatchID: 0x0, Kind: EmatchU32, Flags: 0x0, Pad: 0x0},
						// match 'u32(u16 0x1122 0xffff at nexthdr+4)'
						U32Match: &U32Match{
							Mask:    0xffff,
							Value:   0x1122,
							Off:     0x400,
							OffMask: 0xffff,
						},
					},
				},
			},
			Police: &Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(12)},
		},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalBasic(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Basic{}
			err2 := unmarshalBasic(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Basic missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("marshal(nil)", func(t *testing.T) {
		_, err := marshalBasic(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unmarshal(0x0)", func(t *testing.T) {
		val := Basic{}
		if err := unmarshalBasic([]byte{0x00}, &val); err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
