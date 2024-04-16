package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFq(t *testing.T) {
	weights := []int32{589824, 196608, 65536}
	tests := map[string]struct {
		val  Fq
		err1 error
		err2 error
	}{
		"simple": {
			val: Fq{
				PLimit:           uint32Ptr(1),
				FlowPLimit:       uint32Ptr(2),
				Quantum:          uint32Ptr(3),
				InitQuantum:      uint32Ptr(4),
				RateEnable:       uint32Ptr(5),
				FlowDefaultRate:  uint32Ptr(6),
				FlowMaxRate:      uint32Ptr(7),
				BucketsLog:       uint32Ptr(8),
				FlowRefillDelay:  uint32Ptr(9),
				OrphanMask:       uint32Ptr(10),
				LowRateThreshold: uint32Ptr(11),
				CEThreshold:      uint32Ptr(12),
			},
		},
		"defaults": {
			val: Fq{
				PLimit:           uint32Ptr(10000),
				FlowPLimit:       uint32Ptr(100),
				Quantum:          uint32Ptr(3028),
				InitQuantum:      uint32Ptr(15140),
				RateEnable:       uint32Ptr(1),
				FlowMaxRate:      uint32Ptr(4294967295),
				BucketsLog:       uint32Ptr(10),
				FlowRefillDelay:  uint32Ptr(40000),
				OrphanMask:       uint32Ptr(1023),
				LowRateThreshold: uint32Ptr(68750),
				CEThreshold:      uint32Ptr(4294967295),
				TimerSlack:       uint32Ptr(10000),
				Horizon:          uint32Ptr(10000000),
				HorizonDrop:      uint8Ptr(1),
				PrioMap: &FqPrioQopt{
					Bands:   3,
					PrioMap: [16]uint8{1, 2, 2, 2, 1, 2, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1},
				},
				Weights: &weights,
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFq(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Fq{}
			err2 := unmarshalFq(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Fq missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalFq(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
