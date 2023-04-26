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
					{
						Hdr: EmatchHdr{MatchID: 0x0, Kind: EmatchU32, Flags: EmatchRelAnd, Pad: 0x0},
						U32Match: &U32Match{
							Mask:    0xffff,
							Value:   0x1122,
							Off:     0x4,
							OffMask: 0xffffffff,
						},
					},
					{
						Hdr: EmatchHdr{MatchID: 0x0, Kind: EmatchCmp, Flags: EmatchRelEnd, Pad: 0x0},
						CmpMatch: &CmpMatch{
							Val:   0x14,
							Mask:  0xff00,
							Off:   0x3,
							Align: CmpMatchU16,
							Flags: 0,
							Layer: EmatchLayerTransport,
							Opnd:  EmatchOpndGt,
						},
					},
				},
			},
		},
		"match 'ipset(interactive src,src)'": {
			val: Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1, ProgID: 42},
				Matches: &[]EmatchMatch{
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchIPSet, Flags: EmatchRelEnd, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 19,
							Dir:     []IPSetDir{IPSetSrc, IPSetSrc},
						},
					},
				},
			},
		},
		"match 'ipt(-6 -m foobar)'": {
			val: Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 1, ProgID: 42},
				Matches: &[]EmatchMatch{
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: 0x9, Flags: EmatchRelEnd, Pad: 0x0},
						IptMatch: &IptMatch{
							MatchName: stringPtr("foobar"),
							NFProto:   uint8Ptr(10),
						},
					},
				},
			},
		},

		// A AND (B1 OR B2) AND NOT C
		// EmatchHMatch:
		//   ----------------------------
		//   |  A | (  | C  | B1  | B2)  |
		//   -----------------------------
		"match 'ipset(A src,src)' and  (ipset(B1 dst,dst) or ipset(B2 src,dst,dst)) and not ipset(C src,dst)) ": {
			val: Ematch{
				Hdr: &EmatchTreeHdr{NMatches: 5},
				Matches: &[]EmatchMatch{
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchIPSet, Flags: EmatchRelAnd, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 10, // IpsetName : A
							Dir:     []IPSetDir{IPSetSrc, IPSetSrc},
						},
					},
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchContainer, Flags: EmatchRelAnd, Pad: 0x0},
						ContainerMatch: &ContainerMatch{
							Pos: 3,
						},
					},
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchIPSet, Flags: EmatchInvert, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 13, // IpsetName : C
							Dir:     []IPSetDir{IPSetSrc, IPSetDst},
						},
					},
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchIPSet, Flags: EmatchRelOr, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 11, // IpsetName : B1
							Dir:     []IPSetDir{IPSetDst, IPSetDst},
						},
					},
					{
						Hdr: EmatchHdr{MatchID: 0, Kind: EmatchIPSet, Flags: EmatchRelEnd, Pad: 0x0},
						IPSetMatch: &IPSetMatch{
							IPSetID: 12, // IpsetName : B2
							Dir:     []IPSetDir{IPSetSrc, IPSetDst, IPSetDst},
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
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Ematch{}
			err2 := unmarshalEmatch(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
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
	t.Run("marshal(invalid)", func(t *testing.T) {
		_, err := marshalEmatch(&Ematch{
			Hdr: &EmatchTreeHdr{NMatches: 1, ProgID: 73},
			Matches: &[]EmatchMatch{
				{
					Hdr: EmatchHdr{MatchID: 0, Kind: ematchInvalid},
				},
			},
		})
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
	t.Run("unmarshal(0x0)", func(t *testing.T) {
		val := Ematch{}
		if err := unmarshalEmatch([]byte{0x00}, &val); err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
