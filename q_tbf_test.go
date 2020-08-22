package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTbf(t *testing.T) {
	tests := map[string]struct {
		val  Tbf
		err1 error
		err2 error
	}{
		"no TbfQopt": {val: Tbf{Burst: uint32Ptr(1)}, err1: ErrNoArg},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTbf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Tbf{}
			err2 := unmarshalTbf(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Tbf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalTbf(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
