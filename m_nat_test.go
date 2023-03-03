package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNat(t *testing.T) {
	tests := map[string]struct {
		val  Nat
		err1 error
		err2 error
	}{
		"simple":          {val: Nat{Parms: &NatParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: Nat{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNat(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			newData, tm := injectTcft(t, data, tcaNatTm)
			newData = injectAttribute(t, newData, []byte{}, tcaNatPad)
			val := Nat{}
			err2 := unmarshalNat(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Defact missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalNat(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
