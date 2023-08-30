package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCbs(t *testing.T) {
	tests := map[string]struct {
		val  Cbs
		err1 error
		err2 error
	}{
		"simple": {val: Cbs{Parms: &CbsOpt{Offload: 42}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCbs(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Cbs{}
			err2 := unmarshalCbs(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Cbs missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCbs(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
