package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGact(t *testing.T) {
	tests := map[string]struct {
		val  Gact
		err1 error
		err2 error
	}{
		"failing": {val: Gact{Tm: &Tcft{Install: 2}}, err1: ErrNoArgAlter},
		"simple":  {val: Gact{Parms: &GactParms{Index: 1, Capab: 2}, Prob: &GactProb{PType: 2}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalGact(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(testcase.err1, err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaGactTm)
			newData = injectAttribute(t, newData, []byte{}, tcaGactPad)
			val := Gact{}
			err2 := unmarshalGact(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(testcase.err2, err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Gact missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalGact(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
