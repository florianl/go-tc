package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIpt(t *testing.T) {
	tests := map[string]struct {
		val  Ipt
		err1 error
		err2 error
	}{
		"simple":          {val: Ipt{Table: stringPtr("testTable"), Hook: uint32Ptr(42), Index: uint32Ptr(1984)}},
		"invalidArgument": {val: Ipt{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"simple+Cnt":      {val: Ipt{Table: stringPtr("testTable"), Hook: uint32Ptr(42), Index: uint32Ptr(1984), Cnt: &IptCnt{RefCnt: 7, BindCnt: 42}}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalIpt(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Ipt{}
			err2 := unmarshalIpt(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Ipt missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalIpt(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
