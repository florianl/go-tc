package tc

import (
	"errors"
	"testing"

	"github.com/florianl/go-tc/core"
	"github.com/florianl/go-tc/internal/unix"
)

func TestQdisc(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	err := tcSocket.Qdisc().Add(nil)
	if err != ErrNoArg {
		t.Fatalf("expected ErrNoArg, received: %v", err)
	}

	faultyQdisc := Object{
		Msg: Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: 0,
			Handle:  core.BuildHandle(0xFFFF, 0x0000),
			Parent:  0xFFFFFFF1,
			Info:    0,
		},
	}
	if err := tcSocket.Qdisc().Replace(&faultyQdisc); err != ErrInvalidDev {
		t.Fatalf("expected ErrInvalidDev, received: %v", err)
	}

	tests := map[string]struct {
		kind    string
		err     error
		fqCodel *FqCodel
		red     *Red
		sfb     *Sfb
		sfq     *Sfq
		cbq     *Cbq
		codel   *Codel
		hhf     *Hhf
		pie     *Pie
		choke   *Choke
		netem   *Netem
		cake    *Cake
		htb     *Htb
		prio    *Prio
		plug    *Plug
	}{
		"clsact":   {kind: "clsact"},
		"emptyHtb": {kind: "htb", err: ErrNoArg},
		"fq_codel": {
			kind:    "fq_codel",
			fqCodel: &FqCodel{Target: uint32Ptr(42), Limit: uint32Ptr(0xCAFE)},
		},
		"red": {kind: "red", red: &Red{MaxP: uint32Ptr(42)}},
		"sfb": {kind: "sfb", sfb: &Sfb{Parms: &SfbQopt{Max: 0xFF}}},
		"sfq": {kind: "sfq", sfq: &Sfq{V0: SfqQopt{
			PerturbPeriod: 64,
			Limit:         3000,
			Flows:         512,
		}}},
		"cbq": {kind: "cbq", cbq: &Cbq{
			LssOpt: &CbqLssOpt{OffTime: 10}, WrrOpt: &CbqWrrOpt{Weight: 42},
			FOpt: &CbqFOpt{Split: 2}, OVLStrategy: &CbqOvl{Penalty: 2},
		}},
		"codel": {kind: "codel", codel: &Codel{
			Target: uint32Ptr(1), Limit: uint32Ptr(2), Interval: uint32Ptr(3),
			ECN: uint32Ptr(4), CEThreshold: uint32Ptr(5),
		}},
		"hhf": {kind: "hhf", hhf: &Hhf{
			BacklogLimit: uint32Ptr(1), Quantum: uint32Ptr(2), HHFlowsLimit: uint32Ptr(3),
			ResetTimeout: uint32Ptr(4), AdmitBytes: uint32Ptr(5), EVICTTimeout: uint32Ptr(6), NonHHWeight: uint32Ptr(7),
		}},
		"pie": {kind: "pie", pie: &Pie{
			Target: uint32Ptr(1), Limit: uint32Ptr(2), TUpdate: uint32Ptr(3),
			Alpha: uint32Ptr(4), Beta: uint32Ptr(5), ECN: uint32Ptr(6), Bytemode: uint32Ptr(7),
		}},
		"choke": {kind: "choke", choke: &Choke{MaxP: uint32Ptr(42)}},
		"netem": {kind: "netem", netem: &Netem{Ecn: uint32Ptr(64)}},
		"cake":  {kind: "cake", cake: &Cake{BaseRate: uint64Ptr(128)}},
		"htb":   {kind: "htb", htb: &Htb{Rate64: uint64Ptr(96)}},
		"prio": {kind: "prio", prio: &Prio{
			Bands:   3,
			PrioMap: [16]uint8{1, 2, 2, 2, 1, 2, 9, 9, 1, 1, 1, 1, 1, 1, 1, 1},
		}},
		// TODO(flo): reenable this test.
		//"plug": {kind: "plug", plug: &Plug{Action: PlugReleaseIndefinite}},
	}

	tcMsg := Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: 123,
		Handle:  core.BuildHandle(0xFFFF, 0x0000),
		Parent:  0xFFFFFFF1,
		Info:    0,
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			testQdisc := Object{
				tcMsg,
				Attribute{
					Kind:    testcase.kind,
					FqCodel: testcase.fqCodel,
					Red:     testcase.red,
					Sfb:     testcase.sfb,
					Sfq:     testcase.sfq,
					Cbq:     testcase.cbq,
					Codel:   testcase.codel,
					Hhf:     testcase.hhf,
					Pie:     testcase.pie,
					Choke:   testcase.choke,
					Netem:   testcase.netem,
					Cake:    testcase.cake,
					Htb:     testcase.htb,
					Prio:    testcase.prio,
					Plug:    testcase.plug,
				},
			}

			if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
				if testcase.err != nil && !errors.Is(testcase.err, err) {
					// we received the expected error
					return
				}
				t.Fatalf("could not add new qdisc: %v", err)
			}

			qdiscs, err := tcSocket.Qdisc().Get()
			if err != nil {
				t.Fatalf("could not get qdiscs: %v", err)
			}
			for _, qdisc := range qdiscs {
				t.Logf("%#v\n", qdisc)
			}

			t.Run("Change", func(t *testing.T) {
				if err := tcSocket.Qdisc().Change(&testQdisc); err != nil {
					t.Fatalf("could not change qdisc: %v", err)
				}
			})

			t.Run("Link", func(t *testing.T) {
				if err := tcSocket.Qdisc().Link(&testQdisc); err != nil {
					t.Fatalf("could not change qdisc: %v", err)
				}
			})

			if err := tcSocket.Qdisc().Delete(&testQdisc); err != nil {
				t.Fatalf("could not delete qdisc: %v", err)
			}
		})
	}

	t.Run("delete nil", func(t *testing.T) {
		if err := tcSocket.Qdisc().Delete(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("add nil", func(t *testing.T) {
		if err := tcSocket.Qdisc().Add(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("link nil", func(t *testing.T) {
		if err := tcSocket.Qdisc().Link(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("replace nil", func(t *testing.T) {
		if err := tcSocket.Qdisc().Replace(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("change nil", func(t *testing.T) {
		if err := tcSocket.Qdisc().Change(nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
