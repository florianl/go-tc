package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFqCodel(t *testing.T) {
	tests := map[string]struct {
		val  FqCodel
		err1 error
		err2 error
	}{
		"simple": {val: FqCodel{Target: uint32Ptr(1), Limit: uint32Ptr(2), Interval: uint32Ptr(3), ECN: uint32Ptr(4), Flows: uint32Ptr(5), Quantum: uint32Ptr(6), CEThreshold: uint32Ptr(7), DropBatchSize: uint32Ptr(8), MemoryLimit: uint32Ptr(9)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFqCodel(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := FqCodel{}
			err2 := unmarshalFqCodel(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("FqCodel missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalFqCodel(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
