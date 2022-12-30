package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStab(t *testing.T) {
	foo := []byte{0x13, 0x37}
	tests := map[string]struct {
		val  Stab
		err1 error
		err2 error
	}{
		"simple":  {val: Stab{Base: &SizeSpec{CellLog: 42, LinkLayer: 1}}},
		"simple2": {val: Stab{Base: &SizeSpec{CellLog: 42, LinkLayer: 1}, Data: &foo}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalStab(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Stab{}
			err2 := unmarshalStab(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Stab missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalStab(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
