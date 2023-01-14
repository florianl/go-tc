package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlow(t *testing.T) {
	actions := []*Action{
		{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 0x1, Capab: 0x0, Action: 0x4,
			RefCnt: 0x1, BindCnt: 0x1, Eaction: 0x1, IfIndex: 0x2}}},
	}

	tests := map[string]struct {
		val  Flow
		err1 error
		err2 error
	}{
		"simple": {val: Flow{
			Keys: uint32Ptr(12), Mode: uint32Ptr(34), BaseClass: uint32Ptr(56), RShift: uint32Ptr(78),
			Addend: uint32Ptr(90), Mask: uint32Ptr(21), XOR: uint32Ptr(43), Divisor: uint32Ptr(65), PerTurb: uint32Ptr(87),
		}},
		"with Action": {val: Flow{
			Keys: uint32Ptr(13), Mode: uint32Ptr(31),
			Actions: &actions,
		}},
		"with ematch": {val: Flow{
			Ematch: &Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1},
				Matches: &[]EmatchMatch{
					{
						Hdr: EmatchHdr{MatchID: 0x0, Kind: EmatchU32, Flags: 0x0, Pad: 0x0},
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
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlow(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Flow{}
			err2 := unmarshalFlow(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
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
