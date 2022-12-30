package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenStats(t *testing.T) {
	tests := map[string]struct {
		val  GenStats
		err1 error
		err2 error
	}{
		"simple": {val: GenStats{Basic: &GenBasic{Bytes: 123}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalGenStats(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := GenStats{}
			newData := injectAttribute(t, data, []byte{}, tcaStatsPad)
			err2 := unmarshalGenStats(newData, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("GenStats missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalGenStats(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
