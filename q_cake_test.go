package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCake(t *testing.T) {
	tests := map[string]struct {
		val  Cake
		err1 error
		err2 error
	}{
		"simple": {val: Cake{
			BaseRate:     uint64Ptr(123),
			DiffServMode: uint32Ptr(23),
			Atm:          uint32Ptr(34),
			FlowMode:     uint32Ptr(45),
			Overhead:     uint32Ptr(56),
			Rtt:          uint32Ptr(67),
			Target:       uint32Ptr(78),
			Autorate:     uint32Ptr(89),
			Memory:       uint32Ptr(90),
			Nat:          uint32Ptr(11),
			Raw:          uint32Ptr(22),
			Wash:         uint32Ptr(33),
			Mpu:          uint32Ptr(44),
			Ingress:      uint32Ptr(55),
			AckFilter:    uint32Ptr(66),
			SplitGso:     uint32Ptr(77),
			FwMark:       uint32Ptr(88),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCake(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaCakePad)
			val := Cake{}
			err2 := unmarshalCake(newData, &val)
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
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCake(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCakeXStats(t *testing.T) {
	tests := map[string]struct {
		val  CakeXStats
		err1 error
	}{
		"qdisc": {val: CakeXStats{
			CapacityEstimate: uint64Ptr(12500000),
			MemoryLimit:      uint32Ptr(4194304),
			MemoryUsed:       uint32Ptr(1514),
			AvgNetoff:        uint32Ptr(14),
			MinNetlen:        uint32Ptr(42),
			MaxNetlen:        uint32Ptr(1514),
			MinAdjlen:        uint32Ptr(42),
			MaxAdjlen:        uint32Ptr(1514),
			ActiveQueues:     uint32Ptr(1),
			Tins: []CakeTinStats{
				{
					SentPackets:        uint32Ptr(1),
					SentBytes:          uint64Ptr(2),
					DroppedPackets:     uint32Ptr(3),
					DroppedBytes:       uint64Ptr(4),
					AcksDroppedPackets: uint32Ptr(5),
					AcksDroppedBytes:   uint64Ptr(6),
					EcnMarkedPackets:   uint32Ptr(7),
					EcnMarkedBytes:     uint64Ptr(8),
					BacklogPackets:     uint32Ptr(9),
					BacklogBytes:       uint32Ptr(10),
					ThresholdRate:      uint64Ptr(11),
					TargetUs:           uint32Ptr(12),
					IntervalUs:         uint32Ptr(13),
					WayIndirectHits:    uint32Ptr(14),
					WayMisses:          uint32Ptr(15),
					WayCollisions:      uint32Ptr(16),
					PeakDelayUs:        uint32Ptr(17),
					AvgDelayUs:         uint32Ptr(18),
					BaseDelayUs:        uint32Ptr(19),
					SparseFlows:        uint32Ptr(20),
					BulkFlows:          uint32Ptr(21),
					UnresponsiveFlows:  uint32Ptr(22),
					MaxSkblen:          uint32Ptr(23),
					FlowQuantum:        uint32Ptr(24),
				},
				{SentPackets: uint32Ptr(42), SentBytes: uint64Ptr(4242)},
				{ThresholdRate: uint64Ptr(1000000), TargetUs: uint32Ptr(5000)},
			},
		}},
		"class": {val: CakeXStats{
			Deficit:     int32Ptr(-1514),
			CobaltCount: uint32Ptr(2),
			Dropping:    uint32Ptr(1),
			DropNextUs:  int32Ptr(-42),
			PDrop:       uint32Ptr(73),
			BlueTimerUs: int32Ptr(100),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalCakeXStats(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			newData := injectAttribute(t, data, []byte{}, tcaCakeStatsPad)
			val := CakeXStats{}
			if err := unmarshalCakeXStats(newData, &val); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("CakeXStats missmatch (-want +got):\n%s", diff)
			}

			xs := &XStats{}
			if err := extractXStats(newData, xs, "cake"); err != nil {
				t.Fatalf("extractXStats(cake): %v", err)
			}
			if diff := cmp.Diff(xs.Cake, &testcase.val); diff != "" {
				t.Fatalf("XStats.Cake missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCakeXStats(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unknown attribute", func(t *testing.T) {
		data, err := marshalCakeXStats(&CakeXStats{MemoryUsed: uint32Ptr(42)})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// an attribute type a future kernel might add
		newData := injectAttribute(t, data, []byte{1, 2, 3, 4}, 0x7ff)
		val := CakeXStats{}
		if err := unmarshalCakeXStats(newData, &val); err != nil {
			t.Fatalf("unknown attributes must be skipped: %v", err)
		}
		if diff := cmp.Diff(val, CakeXStats{MemoryUsed: uint32Ptr(42)}); diff != "" {
			t.Fatalf("CakeXStats missmatch (-want +got):\n%s", diff)
		}
	})
}

// TestExtractTcmsgAttributesCakeAppXStats covers the end-to-end attribute
// parsing path for a cake qdisc dump: cake exports its extended statistics
// as nested attributes in TCA_STATS_APP, see cake_dump_stats().
func TestExtractTcmsgAttributesCakeAppXStats(t *testing.T) {
	want := &CakeXStats{
		CapacityEstimate: uint64Ptr(12500000),
		MemoryUsed:       uint32Ptr(1514),
		Tins: []CakeTinStats{
			{SentPackets: uint32Ptr(1), SentBytes: uint64Ptr(2)},
			{SentPackets: uint32Ptr(3), SentBytes: uint64Ptr(4)},
		},
	}
	app, err := marshalXStats(XStats{Cake: want})
	if err != nil {
		t.Fatalf("marshalXStats: %v", err)
	}
	stats2, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStatsApp, Data: app},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}
	attrs, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStats2, Data: stats2},
		{Interpretation: vtString, Type: tcaKind, Data: "cake"},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}

	obj := &Attribute{}
	if err := extractTcmsgAttributes(0, attrs, obj); err != nil {
		t.Fatalf("extractTcmsgAttributes: %v", err)
	}
	if obj.XStats == nil || obj.XStats.Cake == nil {
		t.Fatal("expected XStats.Cake to be populated from TCA_STATS_APP")
	}
	if diff := cmp.Diff(obj.XStats.Cake, want); diff != "" {
		t.Fatalf("XStats.Cake missmatch (-want +got):\n%s", diff)
	}
}
