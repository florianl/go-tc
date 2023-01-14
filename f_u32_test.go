package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestU32(t *testing.T) {
	actions := []*Action{
		{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 0x1, Capab: 0x0, Action: 0x4, RefCnt: 0x1, BindCnt: 0x1, Eaction: 0x1, IfIndex: 0x2}}},
	}

	tests := map[string]struct {
		val  U32
		err1 error
		err2 error
	}{
		"empty": {},
		"simple": {val: U32{
			ClassID: uint32Ptr(0xFFFF),
			Mark:    &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1},
			Hash:    uint32Ptr(1234), Pcnt: uint64Ptr(4321), InDev: stringPtr("foobar"),
		}},
		"divisor": {val: U32{Divisor: uint32Ptr(1), Link: uint32Ptr(42)}},
		"extended": {val: U32{
			ClassID: uint32Ptr(0xFFFF),
			Mark:    &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1},
			Police:  &Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(12)},
		}},
		"policy": {val: U32{
			Sel: &U32Sel{
				Flags: 0x1,
				NKeys: 0x0,
			},
			Police: &Police{
				Tbf: &Policy{
					Action: 0x1,
					Burst:  0xc35000,
					Rate: RateSpec{
						CellLog:   0x3,
						Linklayer: 0x1,
						CellAlign: 0xffff,
						Rate:      0x1e848,
					},
				},
			},
		}},
		"multiple Keys": {val: U32{
			ClassID: uint32Ptr(0xFFFF), Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1},
			Sel: &U32Sel{
				Flags: 0x0,
				NKeys: 0x3,
				Keys: []U32Key{
					{Mask: 0xFF, Val: 0x55, Off: 0x1, OffMask: 0x2},
					{Mask: 0xFF00, Val: 0xAA00, Off: 0x3, OffMask: 0x5},
					{Mask: 0xF0F0, Val: 0x5050, Off: 0xC, OffMask: 0xC},
				},
			},
		}},
		"actions": {val: U32{Flags: uint32Ptr(0x8), Actions: &actions}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalU32(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaU32Pad)
			val := U32{}
			err2 := unmarshalU32(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("U32 missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalU32(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
