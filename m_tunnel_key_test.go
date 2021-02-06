package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTunnelKey(t *testing.T) {
	IPv4 := net.ParseIP("127.0.0.1")
	IPv6 := net.ParseIP("fe80::42")
	tests := map[string]struct {
		val  TunnelKey
		err1 error
		err2 error
	}{
		"simple": {val: TunnelKey{Parms: &TunnelParms{Index: 0x3,
			Capab:   0x0,
			Action:  ActPipe,
			RefCnt:  0x1,
			BindCnt: 0x1, TunnelKeyAction: 0x0},
			KeyEncSrc: &IPv4,
			KeyEncDst: &IPv4}},
		"IPv6": {val: TunnelKey{Parms: &TunnelParms{Index: 42},
			KeyEncSrc: &IPv6, KeyEncDst: &IPv6,
			KeyEncKeyID:   uint32Ptr(0xAA55),
			KeyEncDstPort: uint16Ptr(22),
			KeyNoCSUM:     uint8Ptr(1),
			KeyEncTOS:     uint8Ptr(2),
			KeyEncTTL:     uint8Ptr(42),
		}},
		"invalidArgument": {val: TunnelKey{Tm: &Tcft{Install: 1}},
			err1: ErrNoArgAlter},
	}

	endianessMix := make(map[uint16]valueType)
	endianessMix[tcaTunnelKeyEncKeyID] = vtUint32Be
	endianessMix[tcaTunnelKeyEncDstPort] = vtUint16Be

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTunnelKey(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaTunnelKeyTm)
			newData = injectAttribute(t, newData, []byte{}, tcaTunnelKeyPad)
			newData = changeEndianess(t, newData, endianessMix)
			val := TunnelKey{}
			err2 := unmarshalTunnelKey(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("TunnelKey missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalTunnelKey(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
