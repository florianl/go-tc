package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAtm(t *testing.T) {
	tests := map[string]struct {
		val  Atm
		err1 error
		err2 error
	}{
		"simple": {val: Atm{FD: uint32Ptr(12), Addr: &AtmPvc{Itf: byte(2)}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalAtm(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Atm{}
			err2 := unmarshalAtm(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Atm missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalAtm(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
