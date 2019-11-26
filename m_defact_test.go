package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefact(t *testing.T) {
	var example string = "example"
	tests := map[string]struct {
		val  Defact
		err1 error
		err2 error
	}{
		"simple":          {val: Defact{Parms: &DefactParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: Defact{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"data":            {val: Defact{Data: &example}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalDefact(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Defact{}
			err2 := unmarshalDefact(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Defact missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalDefact(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
