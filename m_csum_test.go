package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCsum(t *testing.T) {
	tests := map[string]struct {
		val  Csum
		err1 error
		err2 error
	}{
		"failing": {val: Csum{Tm: &Tcft{Install: 2}}, err1: ErrNoArgAlter},
		"simple":  {val: Csum{Parms: &CsumParms{Index: 1, Capab: 2}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCsum(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaCsumTm)
			newData = injectAttribute(t, newData, []byte{}, tcaCsumPad)
			val := Csum{}
			err2 := unmarshalCsum(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Csum missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalCsum(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
