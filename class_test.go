package tc

import (
	"testing"

	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
)

func TestClass(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	err := tcSocket.Class().Add(nil)
	if err != ErrNoArg {
		t.Fatalf("expected ErrNoArg, received: %v", err)
	}

	tcMsg := Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: 1337,
		Handle:  core.BuildHandle(0x1, 0x0000),
		Parent:  0xFFFFFFFF,
		Info:    0,
	}

	testQdisc := Object{
		tcMsg,
		Attribute{
			Kind: "htb",
			Htb: &Htb{
				Init: &HtbGlob{
					Defcls: 0x30,
				},
			},
		},
	}

	// tc qdisc add dev $INTERFACE root handle 1: htb default 30
	if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
		t.Fatalf("could not add new qdisc: %v", err)
	}

	tcMsg.Parent = 0x10000
	tcMsg.Handle = core.BuildHandle(0x1, 0x1)

	tests := map[string]struct {
		kind string
		hfsc *Hfsc
	}{
		"hfsc": {kind: "hfsc", hfsc: &Hfsc{Rsc: &ServiceCurve{M1: 12, D: 34, M2: 56}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			testClass := Object{
				tcMsg,
				Attribute{
					Kind: testcase.kind,
					Hfsc: testcase.hfsc,
				},
			}

			if err := tcSocket.Class().Add(&testClass); err != nil {
				t.Fatalf("could not add new class: %v", err)
			}

			classes, err := tcSocket.Class().Get(&tcMsg)
			if err != nil {
				t.Fatalf("could not get classes: %v", err)
			}
			for _, class := range classes {
				t.Logf("%#v\n", class)
			}

			if err := tcSocket.Class().Replace(&testClass); err != nil {
				t.Fatalf("could not replace exisiting class: %v", err)
			}

			if err := tcSocket.Class().Delete(&testClass); err != nil {
				t.Fatalf("could not delete class: %v", err)
			}
		})
	}
}
