package tc

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIfe(t *testing.T) {
	mac, err := net.ParseMAC("00:11:22:33:44:55")
	if err != nil {
		t.Fatalf("failed to parse mac: %v", err)
	}
	tests := map[string]struct {
		val  Ife
		err1 error
		err2 error
	}{
		"empty":           {err1: fmt.Errorf("Ife options are missing")},
		"simple":          {val: Ife{Parms: &IfeParms{Index: 42, Action: 1}}},
		"invalidArgument": {val: Ife{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"macs":            {val: Ife{SMac: &mac, DMac: &mac, Type: uint16Ptr(1)}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalIfe(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaIfeTm)
			newData = injectAttribute(t, newData, []byte{}, tcaIfePad)
			val := Ife{}
			err2 := unmarshalIfe(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Defact missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalIfe(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
