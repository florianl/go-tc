package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestChoke(t *testing.T) {
	tests := map[string]struct {
		val  Choke
		err1 error
		err2 error
	}{
		"simple":   {val: Choke{MaxP: uint32Ptr(42)}},
		"extended": {val: Choke{MaxP: uint32Ptr(43), Parms: &RedQOpt{Limit: 1337}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalChoke(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Choke{}
			err2 := unmarshalChoke(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Choke missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalChoke(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
