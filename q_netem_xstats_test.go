package tc

import (
	"encoding/binary"
	"testing"
)

// TestExtractXStatsNetem verifies that the extended statistics of a netem
// qdisc, as exported via TCA_XSTATS, are parsed instead of aborting the whole
// qdisc dump with "extractXStats(): unsupported kind: netem".
func TestExtractXStatsNetem(t *testing.T) {
	// struct tc_netem_xstats from include/uapi/linux/pkt_sched.h
	data := make([]byte, 7*8)
	for i, v := range []uint64{1, 2, 3, 4, 5, 6, 7} {
		binary.LittleEndian.PutUint64(data[i*8:], v)
	}

	xs := &XStats{}
	if err := extractXStats(data, xs, "netem"); err != nil {
		t.Fatalf("extractXStats(netem): %v", err)
	}
	if xs.Netem == nil {
		t.Fatal("expected Netem xstats to be set")
	}
	if xs.Netem.Delayed != 1 || xs.Netem.Dropped != 2 || xs.Netem.Reordered != 5 ||
		xs.Netem.EcnMarked != 6 || xs.Netem.AllocationErrors != 7 {
		t.Errorf("unexpected netem xstats: %+v", xs.Netem)
	}
}

// TestExtractTcmsgAttributesNetemAppXStats covers the end-to-end attribute
// parsing path for a netem qdisc dump: netem exports its extended statistics
// (struct tc_netem_xstats) via TCA_STATS_APP, see netem_dump_stats() calling
// gnet_stats_copy_app(). The payload must end up in Attribute.XStats.Netem
// instead of aborting the whole dump with "unmarshalGenStats()	4".
func TestExtractTcmsgAttributesNetemAppXStats(t *testing.T) {
	xstatsPayload := make([]byte, 7*8)
	for i, v := range []uint64{1, 2, 3, 4, 5, 6, 7} {
		binary.LittleEndian.PutUint64(xstatsPayload[i*8:], v)
	}

	stats2, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStatsApp, Data: xstatsPayload},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}
	// tcaStats2 comes before tcaKind on purpose: the kernel does not
	// guarantee the attribute order, so the TCA_STATS_APP payload must be
	// parsed after the whole attribute list has been walked.
	attrs, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStats2, Data: stats2},
		{Interpretation: vtString, Type: tcaKind, Data: "netem"},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}

	obj := &Attribute{}
	if err := extractTcmsgAttributes(0, attrs, obj); err != nil {
		t.Fatalf("extractTcmsgAttributes: %v", err)
	}
	if obj.XStats == nil || obj.XStats.Netem == nil {
		t.Fatal("expected XStats.Netem to be populated from TCA_STATS_APP")
	}
	if obj.XStats.Netem.Dropped != 2 || obj.XStats.Netem.Reordered != 5 {
		t.Errorf("unexpected netem xstats: %+v", obj.XStats.Netem)
	}
}

// TestExtractTcmsgAttributesUnknownKindApp ensures that an unknown qdisc kind
// with a TCA_STATS_APP payload is reported as an error.
func TestExtractTcmsgAttributesUnknownKindApp(t *testing.T) {
	payload := make([]byte, 16)
	stats2, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStatsApp, Data: payload},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}
	attrs, err := marshalAttributes([]tcOption{
		{Interpretation: vtBytes, Type: tcaStats2, Data: stats2},
	})
	if err != nil {
		t.Fatalf("marshalAttributes: %v", err)
	}

	obj := &Attribute{Kind: "somenewkind"}
	if err := extractTcmsgAttributes(0, attrs, obj); err == nil {
		t.Fatal("unknown kind with APP payload must fail")
	}
}
