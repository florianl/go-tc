package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHtb(t *testing.T) {
	tests := map[string]struct {
		val  Htb
		err1 error
		err2 error
	}{
		"simple": {val: Htb{Rate64: uint64Ptr(123), Parms: &HtbOpt{Buffer: 0xFFFF}}},
		"extended": {val: Htb{Rate64: uint64Ptr(123), Ceil64: uint64Ptr(321), Parms: &HtbOpt{Buffer: 0xFFFF},
			Offload: boolPtr(true), DirectQlen: uint32Ptr(74), Init: &HtbGlob{DirectPkts: 6789}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHtb(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaHtbPad)
			val := Htb{}
			err2 := unmarshalHtb(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Htb missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalHtb(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
