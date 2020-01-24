package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNetem(t *testing.T) {
	var ecn uint32 = 123
	var lat int64 = 4242
	var jitter int64 = 1337

	tests := map[string]struct {
		val  Netem
		err1 error
		err2 error
	}{
		"simple": {val: Netem{Ecn: &ecn, Latency64: &lat, Jitter64: &jitter}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNetem(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Netem{}
			err2 := unmarshalNetem(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Netem missmatch (want +got):\n%s", diff)
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
