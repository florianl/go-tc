package tc

import (
	"errors"
	"io"
	"testing"

	"github.com/florianl/go-tc/core"
	"github.com/google/go-cmp/cmp"
)

func TestNetem(t *testing.T) {
	delayDist := []int16{9, 7, 5, 3, 1}
	tests := map[string]struct {
		val  Netem
		err1 error
		err2 error
	}{
		"delayDist": {val: Netem{
			Qopt: NetemQopt{
				Latency: core.Time2Tick(2),
				Jitter:  core.Time2Tick(3),
			},
			DelayDist: &delayDist}},
		"simple": {val: Netem{Ecn: uint32Ptr(123), Latency64: int64Ptr(-4242), Jitter64: int64Ptr(-1337)}},
		"qopt":   {val: Netem{Qopt: NetemQopt{Latency: 42}, Rate64: uint64Ptr(1337)}},
		"random": {val: Netem{Corr: &NetemCorr{Delay: 2}, Reorder: &NetemReorder{Correlation: 13}, Corrupt: &NetemCorrupt{Correlation: 11}, Rate: &NetemRate{PacketOverhead: 1337}, Slot: &NetemSlot{MinDelay: 2, MaxDelay: 4}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNetem(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Netem{}
			err2 := unmarshalNetem(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Netem missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil-marshalNetem", func(t *testing.T) {
		_, err := marshalNetem(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("nil-unmarshalNetem", func(t *testing.T) {
		if err := unmarshalNetem([]byte{}, nil); !errors.Is(err, io.EOF) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
