package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTbf(t *testing.T) {
	tests := map[string]struct {
		val  Tbf
		err1 error
		err2 error
	}{
		"empty":    {},
		"simple":   {val: Tbf{Rate64: 1, Prate64: 2, Burst: 3, Pburst: 4}},
		"extended": {val: Tbf{Rate64: 1, Prate64: 2, Burst: 3, Pburst: 4, Parms: &TbfQopt{Buffer: 2, Limit: 3, Mtu: 4}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTbf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Tbf{}
			err2 := unmarshalTbf(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Tbf missmatch (want +got):\n%s", diff)
			}
		})
	}
}
