package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPie(t *testing.T) {
	tests := map[string]struct {
		val  Pie
		err1 error
		err2 error
	}{
		"simple": {val: Pie{Target: 1, Limit: 2, TUpdate: 3, Alpha: 4, Beta: 5, ECN: 6, Bytemode: 7}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPie(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Pie{}
			err2 := unmarshalPie(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Pie missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPie(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
