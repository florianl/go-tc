package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPolice(t *testing.T) {
	tests := map[string]struct {
		val  Police
		err1 error
		err2 error
	}{
		"simple":          {val: Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(42)}},
		"invalidArgument": {val: Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(42), Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"tbfOnly": {val: Police{Tbf: &Policy{
			Index: 0x0, Action: 0x2, Limit: 0x0, Burst: 0x4c4b40, Mtu: 0x2400,
			Rate:     RateSpec{CellLog: 0x6, Linklayer: 0x1, Overhead: 1, CellAlign: 0xffff, Mpu: 1, Rate: 0x7d},
			PeakRate: RateSpec{CellLog: 1, Linklayer: 1, Overhead: 1, CellAlign: 1, Mpu: 1, Rate: 1},
		}}},
		"rate64":     {val: Police{Rate64: uint64Ptr(42)}, err1: ErrNotImplemented},
		"peakrate64": {val: Police{PeakRate64: uint64Ptr(123)}, err1: ErrNotImplemented},
		"rates":      {val: Police{Rate: &RateSpec{Rate: 42}, PeakRate: &RateSpec{Rate: 1337}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPolice(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaPoliceTm)
			newData = injectAttribute(t, newData, []byte{}, tcaPolicePad)
			val := Police{}
			err2 := unmarshalPolice(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Police missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPolice(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
