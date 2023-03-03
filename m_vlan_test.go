package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVLan(t *testing.T) {
	tests := map[string]struct {
		val  VLan
		err1 error
		err2 error
	}{
		"simple":          {val: VLan{Parms: &VLanParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: VLan{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"pushs":           {val: VLan{PushID: uint16Ptr(1), PushProtocol: uint16Ptr(2), PushPriority: uint32Ptr(3)}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalVlan(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaVLanTm)
			newData = injectAttribute(t, newData, []byte{}, tcaVLanPad)
			val := VLan{}
			err2 := unmarshalVLan(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("VLan missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalVlan(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
