package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFq(t *testing.T) {
	tests := map[string]struct {
		val  Fq
		err1 error
		err2 error
	}{
		"empty":  {err1: fmt.Errorf("Fq options are missing")},
		"simple": {val: Fq{PLimit: 1, FlowPLimit: 2, Quantum: 3, InitQuantum: 4, RateEnable: 5, FlowDefaultRate: 6, FlowMaxRate: 7, BucketsLog: 8, FlowRefillDelay: 9, OrphanMask: 10, LowRateThreshold: 11, CEThreshold: 12}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFq(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Fq{}
			err2 := unmarshalFq(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Fq missmatch (want +got):\n%s", diff)
			}
		})
	}
}
