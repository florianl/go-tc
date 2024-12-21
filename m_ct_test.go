package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCt(t *testing.T) {
	tests := map[string]struct {
		val  Ct
		err1 error
		err2 error
	}{
		"simple": {val: Ct{Parms: &CtParms{Index: 3}}},
		"all arguments": {val: Ct{Parms: &CtParms{Capab: 4}, Action: uint16Ptr(5),
			Zone: uint16Ptr(5), Mark: uint32Ptr(0xAA55AA55), MarkMask: uint32Ptr(0x55AA55AA),
			NatIPv4Min: netIPPtr(net.ParseIP("1.2.3.4")), NatIPv4Max: netIPPtr(net.ParseIP("8.8.4.4")),
			NatPortMin: uint16Ptr(42), NatPortMax: uint16Ptr(73),
			HelperName: stringPtr("test"), HelperFamily: uint8Ptr(13), HelperProto: uint8Ptr(14)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCt(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}

			newData, tm := injectTcft(t, data, tcaCtTm)
			newData = injectAttribute(t, newData, []byte{}, tcaCtPad)

			val := Ct{}
			err2 := unmarshalCt(newData, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}

			// Reinject value to expected values
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Ct missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalCt(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
