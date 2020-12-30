package tc

import (
	"errors"
	"testing"

	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
)

func TestChain(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	if err := tcSocket.Chain().Add(nil); err != ErrNoArg {
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

	if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
		t.Fatalf("could not add new qdisc: %v", err)
	}

	tests := map[string]struct {
		chain  uint32
		errAdd error
	}{
		"simple": {chain: 42},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			testChain := Object{
				tcMsg,
				Attribute{
					Chain: &testcase.chain,
				},
			}

			if err := tcSocket.Chain().Add(&testChain); err != nil {
				if testcase.errAdd == nil {
					t.Fatalf("could not add new chain: %v", err)
				}
				// TODO: compare the returned error with the expected one
				return
			}

			chains, err := tcSocket.Chain().Get(&tcMsg)
			if err != nil {
				t.Fatalf("could not get chains: %v", err)
			}
			for _, chain := range chains {
				t.Logf("%#v\n", chain)
			}

			if err := tcSocket.Chain().Delete(&testChain); err != nil {
				t.Fatalf("could not delete chain: %v", err)
			}
		})
	}

	t.Run("delete nil", func(t *testing.T) {
		if err := tcSocket.Chain().Delete(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("add nil", func(t *testing.T) {
		if err := tcSocket.Chain().Add(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("get nil", func(t *testing.T) {
		if _, err := tcSocket.Chain().Get(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
