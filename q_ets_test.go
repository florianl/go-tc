package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEts(t *testing.T) {
	quanta := []uint32{4500, 3000, 2500}
	prioMap := []uint8{0, 1, 1, 1, 2, 3, 4, 5}
	tests := map[string]struct {
		val  Ets
		err1 error
		err2 error
	}{
		// tc qdisc add dev tcDev root handle 1: ets strict 3 quanta 4500 3000 2500 priomap 0 1 1 1 2 3 4 5
		"simple": {val: Ets{NBands: uint8Ptr(6), NStrict: uint8Ptr(3), Quanta: &quanta, PrioMap: &prioMap}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalEts(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Ets{}
			err2 := unmarshalEts(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Ets missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("marshalEts(nil)", func(t *testing.T) {
		_, err := marshalEts(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("marshalEtsQuanta(nil)", func(t *testing.T) {
		_, err := marshalEtsQuanta(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("marshalEtsPrioMap(nil)", func(t *testing.T) {
		_, err := marshalEtsPrioMap(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
