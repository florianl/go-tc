package tc

import (
	"errors"
	"testing"
	"net"

	"github.com/google/go-cmp/cmp"
)

func TestTunnelKey(t *testing.T) {
	var testIP = "127.0.0.1"
	testSrcIP := net.ParseIP(testIP)
	testDstIP := net.ParseIP(testIP)
	tests := map[string]struct {
		val  TunnelKey
		err1 error
		err2 error
	}{
		"simple":          {val: TunnelKey{Parms: &TunnelParms{Index: 0x3,
			Capab:   0x0,
			Action:  ActPipe,
			RefCnt:  0x1,
			BindCnt: 0x1, TunnelKeyAction: 0x0},
			KeyEncIPv4Src: &testSrcIP,
			KeyEncIPv4Dst: &testDstIP}},
		"invalidArgument": {val: TunnelKey{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTunnelKey(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData, tm := injectTcft(t, data, tcaVLanTm)
			newData = injectAttribute(t, newData, []byte{}, tcaTunnelKeyPad)
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
