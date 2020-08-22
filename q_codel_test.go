package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCode(t *testing.T) {
	tests := map[string]struct {
		val  Codel
		err1 error
		err2 error
	}{
		"simple": {val: Codel{Target: uint32Ptr(1), Limit: uint32Ptr(2), Interval: uint32Ptr(3), ECN: uint32Ptr(4), CEThreshold: uint32Ptr(5)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCodel(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Codel{}
			err2 := unmarshalCodel(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Codel missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCodel(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
