package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCtInfo(t *testing.T) {
	tests := map[string]struct {
		val  CtInfo
		err1 error
		err2 error
	}{
		"simple": {val: CtInfo{Act: &CtInfoAct{Action: 13}}},
		"all arguments": {val: CtInfo{Act: &CtInfoAct{RefCnt: 14},
			Zone: uint16Ptr(15), ParmsDscpMask: uint32Ptr(16), ParmsDscpStateMask: uint32Ptr(17),
			ParmsCpMarkMask: uint32Ptr(18), StatsDscpSet: uint64Ptr(19), StatsDscpError: uint64Ptr(20),
			StatsCpMarkSet: uint64Ptr(21)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCtInfo(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := CtInfo{}
			tmp, tm := injectTcft(t, data, tcaCtInfoTm)
			newData := injectAttribute(t, tmp, []byte{}, tcaCtInfoPad)
			err2 := unmarshalCtInfo(newData, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("CtInfo missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalCtInfo(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
