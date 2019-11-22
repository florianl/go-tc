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
		"simple": {val: FqCodel{Target: 1, Limit: 2, Interval: 3, ECN: 4, Flows: 5, Quantum: 6, CEThreshold: 7, DropBatchSize: 8, MemoryLimit: 9}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFqCodel(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := FqCodel{}
			err2 := unmarshalFqCodel(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
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
