package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDrr(t *testing.T) {
	tests := map[string]struct {
		val  Drr
		err1 error
		err2 error
	}{
		"simple": {val: Drr{Quantum: uint32Ptr(234)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalDrr(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Drr{}
			err2 := unmarshalDrr(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Drr missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalDrr(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
