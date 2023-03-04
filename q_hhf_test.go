package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHhf(t *testing.T) {
	tests := map[string]struct {
		val  Hhf
		err1 error
		err2 error
	}{
		"simple": {val: Hhf{BacklogLimit: uint32Ptr(1), Quantum: uint32Ptr(2), HHFlowsLimit: uint32Ptr(3), ResetTimeout: uint32Ptr(4), AdmitBytes: uint32Ptr(5), EVICTTimeout: uint32Ptr(6), NonHHWeight: uint32Ptr(7)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHhf(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Hhf{}
			err2 := unmarshalHhf(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Hhf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalHhf(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
