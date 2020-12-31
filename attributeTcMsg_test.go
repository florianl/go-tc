package tc

import (
	"errors"
	"testing"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/google/go-cmp/cmp"
)

func generatePfifo(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}

	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "pfifo"})
	pfifo, _ := marshalStruct(&FifoOpt{Limit: 123})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: pfifo})

	stats, _ := marshalStruct(&Stats{
		Bytes:      123,
		Packets:    321,
		Drops:      0,
		Overlimits: 42,
	})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStats, Data: stats})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateHtb(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}

	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "htb"})
	htbOption, _ := marshalHtb(&Htb{
		DirectQlen: uint32Ptr(123),
		Rate64:     uint64Ptr(234),
		Ceil64:     uint64Ptr(345),
	})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: htbOption})
	htbXStats, _ := marshalStruct(&HtbXStats{
		Lends:   2,
		Borrows: 3,
		Giants:  4,
		Tokens:  5,
		CTokens: 6,
	})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaXstats, Data: htbXStats})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateClsact(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}

	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "clsact"})
	options = append(options, tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: uint8(96)})
	options = append(options, tcOption{Interpretation: vtUint32, Type: tcaEgressBlock, Data: uint32(4919)})
	options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIngressBlock, Data: uint32(51966)})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: []byte{}})
	options = append(options, tcOption{Interpretation: vtUint32, Type: tcaChain, Data: uint32(42)})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateClsactStab(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "clsact"})
	tmp, _ := marshalStab(&Stab{
		Base: &SizeSpec{
			CellLog:   42,
			LinkLayer: 1,
			MTU:       1492,
		},
	})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStab, Data: tmp})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateMatchall(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "matchall"})
	tmp, _ := marshalMatchall(&Matchall{
		ClassID: uint32Ptr(22),
		Flags:   uint32Ptr(33),
	})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: tmp})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateNetem(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "netem"})
	tmp, _ := marshalNetem(&Netem{Ecn: uint32Ptr(42)})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: tmp})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateCake(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "cake"})
	tmp, _ := marshalCake(&Cake{BaseRate: uint64Ptr(424242)})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: tmp})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func generateQfq(t *testing.T) []byte {
	t.Helper()
	options := []tcOption{}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: "qfq"})
	tmp, _ := marshalQfq(&Qfq{Weight: uint32Ptr(1), Lmax: uint32Ptr(2)})
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: tmp})

	data, err := marshalAttributes(options)
	if err != nil {
		t.Fatalf("could not generate test data: %v", err)
	}
	return data
}

func TestExtractTcmsgAttributes(t *testing.T) {
	tests := map[string]struct {
		input    []byte
		expected *Attribute
		err      error
	}{
		"empty": {input: []byte{}, expected: &Attribute{}},
		"clsact": {input: generateClsact(t), expected: &Attribute{Kind: "clsact", HwOffload: uint8Ptr(0x60),
			EgressBlock: uint32Ptr(0x1337), IngressBlock: uint32Ptr(0xcafe), Chain: uint32Ptr(42)}},
		"htb": {input: generateHtb(t), expected: &Attribute{Kind: "htb",
			XStats: &XStats{Htb: &HtbXStats{Lends: 0x02, Borrows: 0x03, Giants: 0x04, Tokens: 0x05, CTokens: 0x06}},
			Htb:    &Htb{DirectQlen: uint32Ptr(0x7b), Rate64: uint64Ptr(0xea), Ceil64: uint64Ptr(0x0159)}}},
		"pfifo": {input: generatePfifo(t), expected: &Attribute{Kind: "pfifo",
			Pfifo: &FifoOpt{Limit: 123}, Stats: &Stats{Bytes: 123, Packets: 321, Drops: 0, Overlimits: 42}}},
		"clsact+stab": {input: generateClsactStab(t), expected: &Attribute{Kind: "clsact",
			Stab: &Stab{Base: &SizeSpec{CellLog: 0x2a, LinkLayer: 0x01, MTU: 0x05d4}}}},
		"matchall": {input: generateMatchall(t), expected: &Attribute{Kind: "matchall",
			Matchall: &Matchall{ClassID: uint32Ptr(22), Flags: uint32Ptr(33)}}},
		"netem": {input: generateNetem(t), expected: &Attribute{Kind: "netem",
			Netem: &Netem{Ecn: uint32Ptr(42)}}},
		"cake": {input: generateCake(t), expected: &Attribute{Kind: "cake",
			Cake: &Cake{BaseRate: uint64Ptr(424242)}}},
		"qfq": {input: generateQfq(t), expected: &Attribute{Kind: "qfq",
			Qfq: &Qfq{Weight: uint32Ptr(1), Lmax: uint32Ptr(2)}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			value := &Attribute{}
			if err := extractTcmsgAttributes(0xCAFE, testcase.input, value); err != nil {
				if testcase.err != nil && testcase.err.Error() == err.Error() {
					// we received the expected error. everything is fine
					return
				}
				t.Fatalf("Received error '%v', but expected '%v'", err, testcase.err)
			}
			if diff := cmp.Diff(value, testcase.expected); diff != "" {
				t.Fatalf("ExtractTcmsgAttributes missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtractTCAOptions(t *testing.T) {
	tests := map[string]struct {
		kind     string
		data     []byte
		expected *Attribute
		err      error
	}{
		"clsact":         {kind: "clsact", expected: &Attribute{}},
		"clsactWithData": {kind: "clsact", data: []byte{0xde, 0xad, 0xc0, 0xde}, err: ErrInvalidArg},
		"ingress":        {kind: "ingress", expected: &Attribute{}},
		"unknown":        {kind: "unknown", err: ErrUnknownKind},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			value := &Attribute{}
			if err := extractTCAOptions(testcase.data, value, testcase.kind); err != nil {
				if errors.Is(err, testcase.err) {
					// we received the expected error. everything is fine
					return
				}
				t.Fatalf("Received error '%v', but expected '%v'", err, testcase.err)
			}
			if diff := cmp.Diff(value, testcase.expected); diff != "" {
				t.Fatalf("ExtractTcmsgAttributes missmatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestFilterAttribute(t *testing.T) {
	tests := map[string]struct {
		val  *Attribute
		err1 error
		err2 error
	}{
		"basic": {val: &Attribute{Kind: "basic", Basic: &Basic{ClassID: uint32Ptr(2)}}},
		"bpf": {val: &Attribute{Kind: "bpf", BPF: &Bpf{Ops: bytesPtr([]byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}),
			OpsLen:  uint16Ptr(0x1),
			ClassID: uint32Ptr(0x10001),
			Flags:   uint32Ptr(0x1)}}},
		"flow": {val: &Attribute{Kind: "flow", Flow: &Flow{Keys: uint32Ptr(12), Mode: uint32Ptr(34), BaseClass: uint32Ptr(56), RShift: uint32Ptr(78),
			Addend: uint32Ptr(90), Mask: uint32Ptr(21), XOR: uint32Ptr(43), Divisor: uint32Ptr(65), PerTurb: uint32Ptr(87)}}},
		"fw":     {val: &Attribute{Kind: "fw", Fw: &Fw{ClassID: uint32Ptr(12), InDev: stringPtr("lo"), Mask: uint32Ptr(0xFFFF)}}},
		"route4": {val: &Attribute{Kind: "route4", Route4: &Route4{ClassID: uint32Ptr(0xFFFF), To: uint32Ptr(2), From: uint32Ptr(3), IIf: uint32Ptr(4)}}},
		"rsvp":   {val: &Attribute{Kind: "rsvp", Rsvp: &Rsvp{ClassID: uint32Ptr(42), Police: &Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(12)}}}},
		"u32":    {val: &Attribute{Kind: "u32", U32: &U32{ClassID: uint32Ptr(0xFFFF), Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			options, err1 := validateFilterObject(unix.RTM_NEWTFILTER, &Object{Msg{Ifindex: 42}, *testcase.val})
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			data, err := marshalAttributes(options)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			info := &Attribute{}
			err2 := extractTcmsgAttributes(0xCAFE, data, info)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(info, testcase.val); diff != "" {
				t.Fatalf("Filter missmatch (want +got):\n%s", diff)
			}
		})
	}
}

func TestQdiscAttribute(t *testing.T) {
	tests := map[string]struct {
		val  *Attribute
		err1 error
		err2 error
	}{
		"clsact":   {val: &Attribute{Kind: "clsact"}},
		"ingress":  {val: &Attribute{Kind: "ingress"}},
		"atm":      {val: &Attribute{Kind: "atm", Atm: &Atm{FD: uint32Ptr(12), Addr: &AtmPvc{Itf: byte(2)}}}},
		"cbq":      {val: &Attribute{Kind: "cbq", Cbq: &Cbq{LssOpt: &CbqLssOpt{OffTime: 10}, WrrOpt: &CbqWrrOpt{Weight: 42}, FOpt: &CbqFOpt{Split: 2}, OVLStrategy: &CbqOvl{Penalty: 2}}}},
		"codel":    {val: &Attribute{Kind: "codel", Codel: &Codel{Target: uint32Ptr(1), Limit: uint32Ptr(2), Interval: uint32Ptr(3), ECN: uint32Ptr(4), CEThreshold: uint32Ptr(5)}}},
		"drr":      {val: &Attribute{Kind: "drr", Drr: &Drr{Quantum: uint32Ptr(345)}}},
		"dsmark":   {val: &Attribute{Kind: "dsmark", Dsmark: &Dsmark{Indices: uint16Ptr(12), DefaultIndex: uint16Ptr(34), Mask: uint8Ptr(56), Value: uint8Ptr(78)}}},
		"fq":       {val: &Attribute{Kind: "fq", Fq: &Fq{PLimit: uint32Ptr(1), FlowPLimit: uint32Ptr(2), Quantum: uint32Ptr(3), InitQuantum: uint32Ptr(4), RateEnable: uint32Ptr(5), FlowDefaultRate: uint32Ptr(6), FlowMaxRate: uint32Ptr(7), BucketsLog: uint32Ptr(8), FlowRefillDelay: uint32Ptr(9), OrphanMask: uint32Ptr(10), LowRateThreshold: uint32Ptr(11), CEThreshold: uint32Ptr(12)}}},
		"fq_codel": {val: &Attribute{Kind: "fq_codel", FqCodel: &FqCodel{Target: uint32Ptr(1), Limit: uint32Ptr(2), Interval: uint32Ptr(3), ECN: uint32Ptr(4), Flows: uint32Ptr(5), Quantum: uint32Ptr(6), CEThreshold: uint32Ptr(7), DropBatchSize: uint32Ptr(8), MemoryLimit: uint32Ptr(9)}}},
		"hfsc":     {val: &Attribute{Kind: "hfsc", HfscQOpt: &HfscQOpt{DefCls: 42}}},
		"hhf":      {val: &Attribute{Kind: "hhf", Hhf: &Hhf{BacklogLimit: uint32Ptr(1), Quantum: uint32Ptr(2), HHFlowsLimit: uint32Ptr(3), ResetTimeout: uint32Ptr(4), AdmitBytes: uint32Ptr(5), EVICTTimeout: uint32Ptr(6), NonHHWeight: uint32Ptr(7)}}},
		"htb":      {val: &Attribute{Kind: "htb", Htb: &Htb{Init: &HtbGlob{Version: 0x3, Rate2Quantum: 0xa, Defcls: 0x30}}}},
		"mqprio":   {val: &Attribute{Kind: "mqprio", MqPrio: &MqPrio{Mode: uint16Ptr(1), Shaper: uint16Ptr(2), MinRate64: uint64Ptr(3), MaxRate64: uint64Ptr(4)}}},
		"pie":      {val: &Attribute{Kind: "pie", Pie: &Pie{Target: uint32Ptr(1), Limit: uint32Ptr(2), TUpdate: uint32Ptr(3), Alpha: uint32Ptr(4), Beta: uint32Ptr(5), ECN: uint32Ptr(6), Bytemode: uint32Ptr(7)}}},
		"qfq":      {val: &Attribute{Kind: "qfq"}},
		"red":      {val: &Attribute{Kind: "red", Red: &Red{MaxP: uint32Ptr(2), Parms: &RedQOpt{QthMin: 2, QthMax: 4}}}},
		"sfb":      {val: &Attribute{Kind: "sfb", Sfb: &Sfb{Parms: &SfbQopt{Max: 0xFF}}}},
		"tbf":      {val: &Attribute{Kind: "tbf", Tbf: &Tbf{Burst: uint32Ptr(3), Pburst: uint32Ptr(4)}}, err1: ErrNoArg},
		"pfifo":    {val: &Attribute{Kind: "pfifo", Pfifo: &FifoOpt{Limit: 42}}},
		"bfifo":    {val: &Attribute{Kind: "bfifo", Bfifo: &FifoOpt{Limit: 84}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			options, err1 := validateQdiscObject(unix.RTM_NEWQDISC, &Object{Msg{Ifindex: 42}, *testcase.val})
			if err1 != nil {
				if testcase.err1 != nil && errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			data, err := marshalAttributes(options)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			info := &Attribute{}
			err2 := extractTcmsgAttributes(unix.RTM_NEWQDISC, data, info)
			if err2 != nil {
				if testcase.err2 != nil && errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(info, testcase.val); diff != "" {
				t.Fatalf("Filter missmatch (want +got):\n%s", diff)
			}
		})
	}
}

func TestClassAttribute(t *testing.T) {
	tests := map[string]struct {
		val  *Attribute
		err1 error
		err2 error
	}{
		"clsact": {val: &Attribute{Kind: "clsact"}, err1: ErrNotImplemented},
		"hfsc":   {val: &Attribute{Kind: "hfsc", Hfsc: &Hfsc{Rsc: &ServiceCurve{M1: 12, D: 34, M2: 56}}}},
		"qfq":    {val: &Attribute{Kind: "qfq", Qfq: &Qfq{Weight: uint32Ptr(2), Lmax: uint32Ptr(4)}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			options, err1 := validateClassObject(unix.RTM_NEWTCLASS, &Object{Msg{Ifindex: 42}, *testcase.val})
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			data, err := marshalAttributes(options)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			info := &Attribute{}
			err2 := extractTcmsgAttributes(unix.RTM_NEWTCLASS, data, info)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(info, testcase.val); diff != "" {
				t.Fatalf("Filter missmatch (want +got):\n%s", diff)
			}
		})
	}
}
