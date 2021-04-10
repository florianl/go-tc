package tc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGate(t *testing.T) {
	tests := map[string]struct {
		val  Gate
		err1 error
		err2 error
	}{
		"empty": {err1: fmt.Errorf("Gate options are missing")},
		"all options": {val: Gate{Parms: &GateParms{Index: 1}, Priority: int32Ptr(2),
			BaseTime: uint64Ptr(3), CycleTime: uint64Ptr(4), CycleTimeExt: uint64Ptr(5),
			Flags: uint32Ptr(6), ClockID: int32Ptr(-7)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalGate(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaSampleTm)
			newData = injectAttribute(t, newData, []byte{}, tcaGatePad)
			val := Gate{}
			err2 := unmarshalGate(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Gate missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("marshal(nil)", func(t *testing.T) {
		_, err := marshalGate(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unmarshal(0x0)", func(t *testing.T) {
		val := Gate{}
		if err := unmarshalGate([]byte{0x00}, &val); err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
