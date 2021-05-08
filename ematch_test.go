package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmatch(t *testing.T) {

	tests := map[string]struct {
		val  Ematch
		err1 error
		err2 error
	}{
		"empty": {err1: ErrNoArg},
		/*
			"match 'meta(priority eq 0)'": {
				val: Ematch{
					Hdr: &EmatchTreeHdr{NMatches: 1, ProgID: 42},
					Matches: &[]EmatchMatch{
						{Hdr: EmatchHdr{MatchID: 0, Kind: 0x4, Flags: 0x0, Pad: 0x0},
							Data: []byte{0xc, 0x0, 0x1, 0x0, 0x6, 0x10, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x8, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0}},
					},
				},
			},
		*/
		"match 'u32(u16 0x1122 0xffff at nexthdr+4)' and 'cmp(u16 at 3 layer 2 mask 0xff00 gt 20)'": {
			val: Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 2},
				Matches: &[]EmatchMatch{
					{Hdr: EmatchHdr{MatchID: 0x0, Kind: 0x3, Flags: 0x1, Pad: 0x0},
						U32Match: &U32Match{
							Mask:    0xffff,
							Value:   0x1122,
							Off:     0x4,
							OffMask: 0xffffffff,
						}},
					{Hdr: EmatchHdr{MatchID: 0x0, Kind: 0x1, Flags: 0x0, Pad: 0x0},
						CmpMatch: &CmpMatch{
							Val:   0x14,
							Mask:  0xff00,
							Off:   0x3,
							Align: CmpMatchU16,
							Flags: 0,
							Layer: EmatchLayerTransport,
							Opnd:  EmatchOpndGt,
						}},
				},
			},
		},
		"match 'ipset(interactive src,src)'": {
			val: Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1, ProgID: 42},
				Matches: &[]EmatchMatch{
					{Hdr: EmatchHdr{MatchID: 0, Kind: 0x8, Flags: 0x0, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 19,
							Dir:     []IPSetDir{IPSetSrc, IPSetSrc},
						},
					},
				},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalEmatch(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Ematch{}
			err2 := unmarshalEmatch(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Ematch missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("marshal(nil)", func(t *testing.T) {
		_, err := marshalEmatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unmarshal(0x0)", func(t *testing.T) {
		val := Ematch{}
		if err := unmarshalEmatch([]byte{0x00}, &val); err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
