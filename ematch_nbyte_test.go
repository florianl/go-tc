package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNByteMatch(t *testing.T) {
	t.Skip()
	tests := map[string]struct {
		val  NByteMatch
		err1 error
		err2 error
	}{
		"simple": {
			val: NByteMatch{Needle: []byte("helloWorld"),
				Offset: 42,
				Layer:  7},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNByteMatch(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := NByteMatch{}
			err2 := unmarshalNByteMatch(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("NByteMatch missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalNByteMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
