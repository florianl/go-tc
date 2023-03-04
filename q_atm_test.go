package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAtm(t *testing.T) {
	tests := map[string]struct {
		val  Atm
		err1 error
		err2 error
	}{
		"simple": {val: Atm{FD: uint32Ptr(12), Addr: &AtmPvc{Itf: byte(2)}}},
		"extended": {val: Atm{
			FD: uint32Ptr(12), Addr: &AtmPvc{Itf: byte(2)},
			Excess: uint32Ptr(34), State: uint32Ptr(45),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalAtm(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Atm{}
			err2 := unmarshalAtm(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Atm missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalAtm(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
