package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPie(t *testing.T) {
	tests := map[string]struct {
		val  Pie
		err1 error
		err2 error
	}{
		"simple": {val: Pie{Target: uint32Ptr(1), Limit: uint32Ptr(2), TUpdate: uint32Ptr(3), Alpha: uint32Ptr(4), Beta: uint32Ptr(5), ECN: uint32Ptr(6), Bytemode: uint32Ptr(7)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPie(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Pie{}
			err2 := unmarshalPie(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Pie missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPie(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
