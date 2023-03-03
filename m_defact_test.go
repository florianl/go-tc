package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefact(t *testing.T) {
	defactData := "example"
	tests := map[string]struct {
		val  Defact
		err1 error
		err2 error
	}{
		"simple":          {val: Defact{Parms: &DefactParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: Defact{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"data":            {val: Defact{Data: &defactData}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalDefact(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaDefTm)
			newData = injectAttribute(t, newData, []byte{}, tcaDefPad)
			val := Defact{}
			err2 := unmarshalDefact(newData, &val)
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
		_, err := marshalDefact(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
