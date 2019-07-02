package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlow(t *testing.T) {
	tests := map[string]struct {
		val  Flow
		err1 error
		err2 error
	}{
		"empty":  {err1: fmt.Errorf("Flow options are missing")},
		"simple": {val: Flow{Keys: 12, Mode: 34, BaseClass: 56, RShift: 78, Addend: 90, Mask: 21, XOR: 43, Divisor: 65, PerTurb: 87}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlow(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Flow{}
			err2 := unmarshalFlow(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Flow missmatch (want +got):\n%s", diff)
			}
		})
	}
}
