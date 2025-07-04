package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGred(t *testing.T) {
	tests := map[string]struct {
		val  Gred
		err1 error
		err2 error
	}{
		"simple": {val: Gred{
			Parms: &GredQOpt{
				Limit: 42,
			},
			DPS: &GredSOpt{
				DPs: 73,
			},
			MaxP:  uint32Ptr(11),
			Limit: uint32Ptr(42),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalGred(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Gred{}
			err2 := unmarshalGred(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Etf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalGred(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
