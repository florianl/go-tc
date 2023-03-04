package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCake(t *testing.T) {
	tests := map[string]struct {
		val  Cake
		err1 error
		err2 error
	}{
		"simple": {val: Cake{
			BaseRate:     uint64Ptr(123),
			DiffServMode: uint32Ptr(23),
			Atm:          uint32Ptr(34),
			FlowMode:     uint32Ptr(45),
			Overhead:     uint32Ptr(56),
			Rtt:          uint32Ptr(67),
			Target:       uint32Ptr(78),
			Autorate:     uint32Ptr(89),
			Memory:       uint32Ptr(90),
			Nat:          uint32Ptr(11),
			Raw:          uint32Ptr(22),
			Wash:         uint32Ptr(33),
			Mpu:          uint32Ptr(44),
			Ingress:      uint32Ptr(55),
			AckFilter:    uint32Ptr(66),
			SplitGso:     uint32Ptr(77),
			FwMark:       uint32Ptr(88),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCake(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaCakePad)
			val := Cake{}
			err2 := unmarshalCake(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
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
