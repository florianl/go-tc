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
		"with ematch": {val: Basic{ClassID: uint32Ptr(3),
			Ematch: &Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1},
				Matches: &[]EmatchMatch{
					{Hdr: EmatchHdr{MatchID: 0x0, Kind: 0x1, Flags: 0x0, Pad: 0x0},
						Data: []byte{0x14, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x3, 0x0, 0x2, 0x12}}},
			}}},
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
