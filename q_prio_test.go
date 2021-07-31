package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrio(t *testing.T) {
	tests := map[string]struct {
		val  Prio
		err1 error
		err2 error
	}{
		"pfifo_fast": {val: Prio{Bands: 3,
			PrioMap: [16]uint8{1, 2, 2, 2, 1, 2, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPrio(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Prio{}
			err2 := unmarshalPrio(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Prio missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPrio(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
