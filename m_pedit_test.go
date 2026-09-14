package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPedit(t *testing.T) {
	tests := map[string]struct {
		val  Pedit
		err1 error
		err2 error
	}{
		"simple": {val: Pedit{
			Sel:    PeditSel{Action: 1},
			Keys:   []PeditKey{{Val: 0xc0a80101, Off: 12}},
			KeysEx: []PeditKeyEx{{HeaderType: PeditHeaderTypeIP4, Cmd: PeditCmdSet}},
		}},
		"invalidArgument":   {val: Pedit{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"keyLengthMismatch": {val: Pedit{Keys: []PeditKey{{Val: 1}}, KeysEx: []PeditKeyEx{}}, err1: errors.New("pedit keys and extended keys length mismatch")},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalPedit(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaPeditPad)
			val := Pedit{}
			err2 := unmarshalPedit(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			// marshalPedit normalizes NKeys from len(Keys); reflect that in expected value.
			testcase.val.Sel.NKeys = uint8(len(testcase.val.Keys))
			if diff := cmp.Diff(val.Sel, testcase.val.Sel); diff != "" {
				t.Fatalf("PeditSel mismatch (want +got):\n%s", diff)
			}
			if diff := cmp.Diff(val.Keys, testcase.val.Keys); diff != "" {
				t.Fatalf("PeditKeys mismatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalPedit(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unmarshalPedit()", func(t *testing.T) {
		err := unmarshalPedit([]byte{0x0}, nil)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}

func TestPeditHelpers(t *testing.T) {
	t.Run("SetIPv4Src", func(t *testing.T) {
		p := &Pedit{}
		ip := net.ParseIP("192.168.1.1")
		p.SetIPv4Src(ip)
		if len(p.Keys) != 1 || len(p.KeysEx) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.KeysEx[0].HeaderType != PeditHeaderTypeIP4 {
			t.Fatalf("expected IP4 header type, got %v", p.KeysEx[0].HeaderType)
		}
		if p.KeysEx[0].Cmd != PeditCmdSet {
			t.Fatalf("expected Set cmd, got %v", p.KeysEx[0].Cmd)
		}
		if p.Keys[0].Off != 12 {
			t.Fatalf("expected offset 12 (src IP), got %d", p.Keys[0].Off)
		}
	})

	t.Run("SetIPv4Dst", func(t *testing.T) {
		p := &Pedit{}
		ip := net.ParseIP("10.0.0.1")
		p.SetIPv4Dst(ip)
		if len(p.Keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.Keys[0].Off != 16 {
			t.Fatalf("expected offset 16 (dst IP), got %d", p.Keys[0].Off)
		}
	})

	t.Run("SetIPv4Src_nonIPv4", func(t *testing.T) {
		p := &Pedit{}
		p.SetIPv4Src(net.ParseIP("::1"))
		if len(p.Keys) != 0 {
			t.Fatalf("expected no keys for IPv6 address, got %d", len(p.Keys))
		}
	})

	t.Run("SetUDPSrc", func(t *testing.T) {
		p := &Pedit{}
		p.SetUDPSrc(8080)
		if len(p.Keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.KeysEx[0].HeaderType != PeditHeaderTypeUDP {
			t.Fatalf("expected UDP header type, got %v", p.KeysEx[0].HeaderType)
		}
	})

	t.Run("SetUDPDst", func(t *testing.T) {
		p := &Pedit{}
		p.SetUDPDst(443)
		if len(p.Keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.KeysEx[0].HeaderType != PeditHeaderTypeUDP {
			t.Fatalf("expected UDP header type, got %v", p.KeysEx[0].HeaderType)
		}
	})

	t.Run("SetSrcPort_TCP", func(t *testing.T) {
		p := &Pedit{}
		p.SetSrcPort(1234, peditIPProtoTCP)
		if len(p.Keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.KeysEx[0].HeaderType != PeditHeaderTypeTCP {
			t.Fatalf("expected TCP header type, got %v", p.KeysEx[0].HeaderType)
		}
	})

	t.Run("SetSrcPort_unknownProtocol", func(t *testing.T) {
		p := &Pedit{}
		p.SetSrcPort(1234, 253)
		if len(p.Keys) != 0 {
			t.Fatalf("expected no keys for unknown protocol, got %d", len(p.Keys))
		}
	})

	t.Run("SetEthDst", func(t *testing.T) {
		p := &Pedit{}
		mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
		p.SetEthDst(mac)
		if len(p.Keys) != 2 {
			t.Fatalf("expected 2 keys for MAC dst, got %d", len(p.Keys))
		}
		if p.KeysEx[0].HeaderType != PeditHeaderTypeEth {
			t.Fatalf("expected Eth header type, got %v", p.KeysEx[0].HeaderType)
		}
	})

	t.Run("SetEthSrc", func(t *testing.T) {
		p := &Pedit{}
		mac, _ := net.ParseMAC("11:22:33:44:55:66")
		p.SetEthSrc(mac)
		if len(p.Keys) != 2 {
			t.Fatalf("expected 2 keys for MAC src, got %d", len(p.Keys))
		}
	})

	t.Run("SetEthDst_shortMAC", func(t *testing.T) {
		p := &Pedit{}
		p.SetEthDst(net.HardwareAddr{0xaa, 0xbb})
		if len(p.Keys) != 0 {
			t.Fatalf("expected no keys for short MAC, got %d", len(p.Keys))
		}
	})

	t.Run("AddIPv4TTL", func(t *testing.T) {
		p := &Pedit{}
		p.AddIPv4TTL(1)
		if len(p.Keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(p.Keys))
		}
		if p.KeysEx[0].Cmd != PeditCmdAdd {
			t.Fatalf("expected Add cmd, got %v", p.KeysEx[0].Cmd)
		}
		if p.Keys[0].Off != 8 {
			t.Fatalf("expected offset 8 (TTL), got %d", p.Keys[0].Off)
		}
	})

	t.Run("NKeys_tracked", func(t *testing.T) {
		p := &Pedit{}
		p.SetIPv4Src(net.ParseIP("1.2.3.4"))
		p.SetIPv4Dst(net.ParseIP("5.6.7.8"))
		if p.Sel.NKeys != 2 {
			t.Fatalf("expected NKeys=2, got %d", p.Sel.NKeys)
		}
	})
}
