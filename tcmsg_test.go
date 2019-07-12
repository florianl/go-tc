//+build linux

package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sys/unix"
)

func TestMsg(t *testing.T) {
	tests := map[string]struct {
		Family    uint32
		IfIndex   uint32
		Handle    uint32
		Parent    uint32
		Info      uint32
		encodeErr error
		decodeErr error
	}{
		"AF_UNSPEC": {Family: unix.AF_UNSPEC, IfIndex: 0, Handle: 0xFFAAFFAA, Parent: Ingress, Info: 0},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := tcmsgEncode(&Msg{
				Family:  testcase.Family,
				Ifindex: testcase.IfIndex,
				Handle:  testcase.Handle,
				Parent:  testcase.Parent,
				Info:    testcase.Info,
			})
			if err != nil {
				if testcase.encodeErr == nil {
					t.Fatalf("expected no encoding error, but got: %v", err)
				}
				// TODO compare resulting errors
				return
			}
			var msg Msg
			if err := tcmsgDecode(data, &msg); err != nil {
				if testcase.decodeErr == nil {

					t.Fatalf("expected no decoding error, but got: %v", err)
				}
				// TODO compare resulting errors
				return
			}
			if diff := cmp.Diff(Msg{
				Family:  testcase.Family,
				Ifindex: testcase.IfIndex,
				Handle:  testcase.Handle,
				Parent:  testcase.Parent,
				Info:    testcase.Info,
			}, msg); diff != "" {
				t.Fatalf("unexpected Msg value (-want +got):\n%s", diff)
			}

		})
	}
}
