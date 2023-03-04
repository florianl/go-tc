package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDsmark(t *testing.T) {
	tests := map[string]struct {
		val  Dsmark
		err1 error
		err2 error
	}{
		"simple":         {val: Dsmark{Indices: uint16Ptr(12), DefaultIndex: uint16Ptr(34), Mask: uint8Ptr(56), Value: uint8Ptr(78)}},
		"simpleWithFlag": {val: Dsmark{Indices: uint16Ptr(12), DefaultIndex: uint16Ptr(34), SetTCIndex: boolPtr(true)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalDsmark(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Dsmark{}
			err2 := unmarshalDsmark(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Dsmark missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalDsmark(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
