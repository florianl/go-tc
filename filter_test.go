//+build linux

package tc

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestFilter(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	err := tcSocket.Filter().Add(nil)
	if err != ErrNoArg {
		t.Fatalf("expected ErrNoArg, received: %v", err)
	}

	tcMsg := Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: 1337,
		Handle:  BuildHandle(0xFFFF, 0x0000),
		Parent:  0xFFFFFFF1,
		Info:    0,
	}

	testQdisc := Object{
		tcMsg,
		Attribute{
			Kind: "clsact",
		},
	}

	u32ExactMatch := &U32{}

	if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
		t.Fatalf("could not add new qdisc: %v", err)
	}

	tests := map[string]struct {
		kind string
		u32  *U32
	}{
		"u32-exactMatch": {kind: "u32", u32: u32ExactMatch},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {

			testFilter := Object{
				tcMsg,
				Attribute{
					Kind: testcase.kind,
					U32:  testcase.u32,
				},
			}

			if err := tcSocket.Filter().Add(&testFilter); err != nil {
				t.Fatalf("could not add new filter: %v", err)
			}

			filters, err := tcSocket.Filter().Get(&tcMsg)
			if err != nil {
				t.Fatalf("could not get filters: %v", err)
			}
			for _, filter := range filters {
				t.Logf("%#v\n", filter)
			}
			if err := tcSocket.Filter().Delete(&testFilter); err != nil {
				t.Fatalf("could not delete filter: %v", err)
			}
		})
	}

}
