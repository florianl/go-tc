package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMPLS(t *testing.T) {
	tests := map[string]struct {
		val  MPLS
		err1 error
		err2 error
	}{
		"all options": {val: MPLS{
			Parms: &MPLSParam{
				Index:   1,
				MAction: MPLSActModify,
			},
			Proto: int16Ptr(101),
			Label: uint32Ptr(102),
			TC:    uint8Ptr(103),
			TTL:   uint8Ptr(104),
			BOS:   uint8Ptr(105),
		}},
		"tm": {
			val: MPLS{
				Tm: &Tcft{
					Install: 1,
					LastUse: 2,
				},
			},
			err1: ErrNoArgAlter,
		},
	}

	endianessMix := make(map[uint16]valueType)
	endianessMix[tcaMPLSProto] = vtInt16Be

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMPLS(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaMPLSTm)
			newData = changeEndianess(t, newData, endianessMix)
			val := MPLS{}
			err2 := unmarshalMPLS(newData, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("MPLS missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("marshal(nil)", func(t *testing.T) {
		_, err := marshalMPLS(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unmarshal(0x0)", func(t *testing.T) {
		val := MPLS{}
		if err := unmarshalMPLS([]byte{0x00}, &val); err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
