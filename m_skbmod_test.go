package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSkbMod(t *testing.T) {
	srcMac, _ := net.ParseMAC("00:00:5e:00:53:01")
	dstMac, _ := net.ParseMAC("00:00:5e:00:53:02")
	tests := map[string]struct {
		val  SkbMod
		err1 error
		err2 error
	}{
		"simple": {val: SkbMod{
			Parms: &SkbModParms{Index: 42},
			SMac:  &srcMac,
			DMac:  &dstMac,
			EType: uint16Ptr(13)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalSkbMod(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			newData, tm := injectTcft(t, data, tcaSkbModTm)
			newData = injectAttribute(t, newData, []byte{}, tcaSkbModPad)
			val := SkbMod{}
			err2 := unmarshalSkbMod(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
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
		_, err := marshalSkbMod(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unmarshalSkbMod()", func(t *testing.T) {
		err := unmarshalSkbMod([]byte{0x0}, nil)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}
