package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCbq(t *testing.T) {
	tests := map[string]struct {
		val  Cbq
		err1 error
		err2 error
	}{
		"empty":  {},
		"simple": {val: Cbq{LssOpt: &CbqLssOpt{OffTime: 10}, WrrOpt: &CbqWrrOpt{Weight: 42}, FOpt: &CbqFOpt{Split: 2}, OVLStrategy: &CbqOvl{Penalty: 2}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCbq(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Cbq{}
			err2 := unmarshalCbq(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Cbq missmatch (want +got):\n%s", diff)
			}
		})
	}
}
