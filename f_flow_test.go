package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlow(t *testing.T) {
	tests := map[string]struct {
		val  Flow
		err1 error
		err2 error
	}{
		"simple": {val: Flow{Keys: uint32Ptr(12), Mode: uint32Ptr(34), BaseClass: uint32Ptr(56), RShift: uint32Ptr(78),
			Addend: uint32Ptr(90), Mask: uint32Ptr(21), XOR: uint32Ptr(43), Divisor: uint32Ptr(65), PerTurb: uint32Ptr(87)}},
		"with ematch": {val: Flow{
			Ematch: &Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1},
				Matches: &[]EmatchMatch{
					{Hdr: EmatchHdr{MatchID: 0x0, Kind: 0x1, Flags: 0x0, Pad: 0x0},
						Data: []byte{0x14, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x3, 0x0, 0x2, 0x12}}},
			}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlow(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Flow{}
			err2 := unmarshalFlow(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Flow missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalFlow(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
