package tc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSample(t *testing.T) {
	tests := map[string]struct {
		val  Sample
		err1 error
		err2 error
	}{
		"empty": {err1: fmt.Errorf("Sample options are missing")},
		"simple": {val: Sample{
			Parms: &SampleParms{Index: 42, Action: 1},
			Rate:  uint32Ptr(42), TruncSize: uint32Ptr(1337), SampleGroup: uint32Ptr(11),
		}},
		"invalidArgument": {val: Sample{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalSample(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaSampleTm)
			newData = injectAttribute(t, newData, []byte{}, tcaSamplePad)
			val := Sample{}
			err2 := unmarshalSample(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Sample missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalSample(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
