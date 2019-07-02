package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDsmark(t *testing.T) {
	tests := map[string]struct {
		val  Dsmark
		err1 error
		err2 error
	}{
		"empty":  {err1: fmt.Errorf("Dsmark options are missing")},
		"simple": {val: Dsmark{Indices: 12, DefaultIndex: 34, Mask: 56, Value: 78}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalDsmark(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Dsmark{}
			err2 := unmarshalDsmark(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Dsmark missmatch (want +got):\n%s", diff)
			}
		})
	}
}
