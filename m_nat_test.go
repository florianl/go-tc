package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNat(t *testing.T) {
	tests := map[string]struct {
		val  Nat
		err1 error
		err2 error
	}{
		"empty":           {err1: fmt.Errorf("Nat options are missing")},
		"simple":          {val: Nat{Parms: &NatParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: Nat{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNat(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Nat{}
			err2 := unmarshalNat(data, &val)
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
}
