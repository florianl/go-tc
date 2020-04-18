package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCake(t *testing.T) {
	var baseRate uint64 = 123
	var diffServMode, atm, flowMode, overhead, rtt, target uint32 = 23, 34, 45, 56, 67, 78
	var autorate, memory, nat, raw, wash, mpu, ingress, ackFilter, splitGso, fwMark uint32 = 89, 90, 11, 22, 33, 44, 55, 66, 77, 88
	tests := map[string]struct {
		val  Cake
		err1 error
		err2 error
	}{
		"simple": {val: Cake{BaseRate: &baseRate,
			DiffServMode: &diffServMode,
			Atm:          &atm,
			FlowMode:     &flowMode,
			Overhead:     &overhead,
			Rtt:          &rtt,
			Target:       &target,
			Autorate:     &autorate,
			Memory:       &memory,
			Nat:          &nat,
			Raw:          &raw,
			Wash:         &wash,
			Mpu:          &mpu,
			Ingress:      &ingress,
			AckFilter:    &ackFilter,
			SplitGso:     &splitGso,
			FwMark:       &fwMark}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCake(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Cake{}
			err2 := unmarshalCake(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Netem missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCake(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
