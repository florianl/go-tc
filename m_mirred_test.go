package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMirred(t *testing.T) {
	tests := map[string]struct {
		val  Mirred
		err1 error
		err2 error
	}{
		"empty":           {err1: fmt.Errorf("Mirred options are missing")},
		"simple":          {val: Mirred{Parms: &MirredParam{Index: 42, Action: 1}}},
		"invalidArgument": {val: Mirred{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMirred(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Mirred{}
			err2 := unmarshalMirred(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Mirred missmatch (want +got):\n%s", diff)
			}
		})
	}
}
