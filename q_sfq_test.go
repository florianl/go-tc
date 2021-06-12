package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSfq(t *testing.T) {
	tests := map[string]struct {
		val  Sfq
		err1 error
		err2 error
	}{
		"simple": {val: Sfq{V0: SfqQopt{
			PerturbPeriod: 64,
			Limit:         3000,
			Flows:         512,
		}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalSfq(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Sfq{}
			err2 := unmarshalSfq(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Sfq missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalSfq(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
