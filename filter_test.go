package tc

import (
	"testing"

	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
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
		Handle:  core.BuildHandle(0xFFFF, 0x0000),
		Parent:  0xFFFFFFF1,
		Info:    0,
	}

	testQdisc := Object{
		tcMsg,
		Attribute{
			Kind: "clsact",
		},
	}

	u32ExactMatch := &U32{
		ClassID: uint32Ptr(13),
	}

	if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
		t.Fatalf("could not add new qdisc: %v", err)
	}

	tests := map[string]struct {
		kind       string
		u32        *U32
		errAdd     error
		errReplace error
	}{
		"unknown":         {kind: "unknown", errAdd: ErrNotImplemented},
		"missingArgument": {kind: "bpf", errAdd: ErrNoArg},
		"u32-exactMatch":  {kind: "u32", u32: u32ExactMatch},
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
				if testcase.errAdd == nil {
					t.Fatalf("could not add new filter: %v", err)
				}
				// TODO: compare the returned error with the expected one
				return
			}

			filters, err := tcSocket.Filter().Get(&tcMsg)
			if err != nil {
				t.Fatalf("could not get filters: %v", err)
			}
			for _, filter := range filters {
				t.Logf("%#v\n", filter)
			}

			if err := tcSocket.Filter().Replace(&testFilter); err != nil {
				if testcase.errReplace == nil {
					t.Fatalf("could not replace filter: %v", err)
				}
				// TODO: compare the returned error with the expected one
				return
			}

			if err := tcSocket.Filter().Delete(&testFilter); err != nil {
				t.Fatalf("could not delete filter: %v", err)
			}
		})
	}

}
