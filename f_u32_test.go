package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestU32(t *testing.T) {
	tests := map[string]struct {
		val  U32
		err1 error
		err2 error
	}{
		"empty":    {},
		"simple":   {val: U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}}},
		"extended": {val: U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}, Police: &Police{AvRate: 1337, Result: 12}}},
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
		"multiple Keys": {val: U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1},
			Sel: &U32Sel{
				Flags: 0x0,
				NKeys: 0x3,
				Keys: []U32Key{
					{Mask: 0xFF, Val: 0x55, Off: 0x1, OffMask: 0x2},
					{Mask: 0xFF00, Val: 0xAA00, Off: 0x3, OffMask: 0x5},
					{Mask: 0xF0F0, Val: 0x5050, Off: 0xC, OffMask: 0xC},
				},
			}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalU32(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := U32{}
			err2 := unmarshalU32(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("U32 missmatch (want +got):\n%s", diff)
			}
		})
	}
}
