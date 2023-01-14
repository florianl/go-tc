package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTcIndex(t *testing.T) {
	actions := []*Action{
		{Kind: "csum", CSum: &Csum{Parms: &CsumParms{Index: 4, Capab: 5}}},
	}

	tests := map[string]struct {
		val  TcIndex
		err1 error
		err2 error
	}{
		"empty":       {},
		"with Action": {val: TcIndex{Hash: uint32Ptr(0xAA55AA55), Actions: &actions}},
		"simple": {val: TcIndex{
			Hash: uint32Ptr(1), Mask: uint16Ptr(2), Shift: uint32Ptr(3),
			FallThrough: uint32Ptr(4), ClassID: uint32Ptr(5),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTcIndex(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := TcIndex{}
			err2 := unmarshalTcIndex(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("TcIndex missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalTcIndex(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
