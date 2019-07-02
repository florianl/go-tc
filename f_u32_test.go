package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestU32(t *testing.T) {
	tests := map[string]struct {
		val  U32
		err1 error
		err2 error
	}{
		"empty":    {err1: fmt.Errorf("U32 options are missing")},
		"simple":   {val: U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}}},
		"extended": {val: U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}, Police: &Police{AvRate: 1337, Result: 12}}},
		"policy": {val: U32{
			Sel: &U32Sel{
				Flags: 0x1,
				NKeys: 0x1,
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
