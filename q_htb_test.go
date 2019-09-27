package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHtb(t *testing.T) {
	tests := map[string]struct {
		val  Htb
		err1 error
		err2 error
	}{
		"empty":    {},
		"simple":   {val: Htb{Rate64: 123, Parms: &HtbOpt{Buffer: 0xFFFF}}},
		"extended": {val: Htb{Rate64: 123, Parms: &HtbOpt{Buffer: 0xFFFF}, Init: &HtbGlob{DirectPkts: 6789}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHtb(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Htb{}
			err2 := unmarshalHtb(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Htb missmatch (want +got):\n%s", diff)
			}
		})
	}
}
