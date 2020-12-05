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
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Stab{}
			err2 := unmarshalStab(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Stab missmatch (-want +got):\n%s", diff)
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
