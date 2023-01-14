package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFw(t *testing.T) {
	actions := []*Action{
		{Kind: "gate", Gate: &Gate{
			Parms: &GateParms{Index: 1}, Priority: int32Ptr(2),
			BaseTime: uint64Ptr(3), CycleTime: uint64Ptr(4), CycleTimeExt: uint64Ptr(5),
			Flags: uint32Ptr(6), ClockID: int32Ptr(-7),
		}},
	}

	tests := map[string]struct {
		val  Fw
		err1 error
		err2 error
	}{
		"simple":   {val: Fw{ClassID: uint32Ptr(12), InDev: stringPtr("lo"), Mask: uint32Ptr(0xFFFF)}},
		"extended": {val: Fw{ClassID: uint32Ptr(12), InDev: stringPtr("lo"), Mask: uint32Ptr(0xFFFF), Police: &Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(12)}}},
		"mixed":    {val: Fw{ClassID: uint32Ptr(12), InDev: stringPtr("lo"), Actions: &actions}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFw(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Fw{}
			err2 := unmarshalFw(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Fw missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalFw(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
