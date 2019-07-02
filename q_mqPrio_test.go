package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMqPrio(t *testing.T) {
	tests := map[string]struct {
		val  MqPrio
		err1 error
		err2 error
	}{
		"empty":  {err1: fmt.Errorf("MqPrio options are missing")},
		"simple": {val: MqPrio{Mode: 1, Shaper: 2, MinRate64: 3, MaxRate64: 4}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMqPrio(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := MqPrio{}
			err2 := unmarshalMqPrio(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("MqPrio missmatch (want +got):\n%s", diff)
			}
		})
	}
}
