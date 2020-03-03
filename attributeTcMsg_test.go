package tc

import (
	"strings"
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
		DirectQlen: 123,
		Rate64:     234,
		Ceil64:     345,
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

func TestExtractTcmsgAttributes(t *testing.T) {
	tests := map[string]struct {
		input    []byte
		expected *Attribute
		err      error
	}{
		"empty":  {input: []byte{}, expected: &Attribute{}},
		"clsact": {input: generateClsact(t), expected: &Attribute{Kind: "clsact", HwOffload: 0x60, EgressBlock: 0x1337, IngressBlock: 0xcafe, Chain: 42}},
		"htb": {input: generateHtb(t), expected: &Attribute{Kind: "htb",
			XStats: &XStats{Htb: &HtbXStats{Lends: 0x02, Borrows: 0x03, Giants: 0x04, Tokens: 0x05, CTokens: 0x06}},
			Htb:    &Htb{DirectQlen: 0x7b, Rate64: 0xea, Ceil64: 0x0159}}},
		"pfifo": {input: generatePfifo(t), expected: &Attribute{Kind: "pfifo",
			Pfifo: &FifoOpt{Limit: 123}, Stats: &Stats{Bytes: 123, Packets: 321, Drops: 0, Overlimits: 42}}},
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
		err      string
	}{
		"clsact":         {kind: "clsact", expected: &Attribute{}},
		"clsactWithData": {kind: "clsact", data: []byte{0xde, 0xad, 0xc0, 0xde}, expected: &Attribute{}, err: "extractClsact()"},
		"ingress":        {kind: "ingress", expected: &Attribute{}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			value := &Attribute{}
			if err := extractTCAOptions(testcase.data, value, testcase.kind); err != nil {
				if len(testcase.err) > 0 && strings.Contains(err.Error(), testcase.err) {
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
		"basic": {val: &Attribute{Kind: "basic", Basic: &Basic{ClassID: 2}}},
		"bpf": {val: &Attribute{Kind: "bpf", BPF: &Bpf{Ops: []byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff},
			OpsLen:  0x1,
			ClassID: 0x10001,
			Flags:   0x1}}},
		"flow":   {val: &Attribute{Kind: "flow", Flow: &Flow{Keys: 12, Mode: 34, BaseClass: 56, RShift: 78, Addend: 90, Mask: 21, XOR: 43, Divisor: 65, PerTurb: 87}}},
		"fw":     {val: &Attribute{Kind: "fw", Fw: &Fw{ClassID: 12, InDev: "lo", Mask: 0xFFFF}}},
		"route4": {val: &Attribute{Kind: "route4", Route4: &Route4{ClassID: 0xFFFF, To: 2, From: 3, IIf: 4}}},
		"rsvp":   {val: &Attribute{Kind: "rsvp", Rsvp: &Rsvp{ClassID: 42, Police: &Police{AvRate: 1337, Result: 12}}}},
		"u32":    {val: &Attribute{Kind: "u32", U32: &U32{ClassID: 0xFFFF, Mark: &U32Mark{Val: 0x55, Mask: 0xAA, Success: 0x1}}}},
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
		"atm":      {val: &Attribute{Kind: "atm", Atm: &Atm{FD: 12, Addr: &AtmPvc{Itf: byte(2)}}}},
		"cbq":      {val: &Attribute{Kind: "cbq", Cbq: &Cbq{LssOpt: &CbqLssOpt{OffTime: 10}, WrrOpt: &CbqWrrOpt{Weight: 42}, FOpt: &CbqFOpt{Split: 2}, OVLStrategy: &CbqOvl{Penalty: 2}}}},
		"codel":    {val: &Attribute{Kind: "codel", Codel: &Codel{Target: 1, Limit: 2, Interval: 3, ECN: 4, CEThreshold: 5}}},
		"drr":      {val: &Attribute{Kind: "drr", Drr: &Drr{Quantum: 345}}},
		"dsmark":   {val: &Attribute{Kind: "dsmark", Dsmark: &Dsmark{Indices: 12, DefaultIndex: 34, Mask: 56, Value: 78}}},
		"fq":       {val: &Attribute{Kind: "fq", Fq: &Fq{PLimit: 1, FlowPLimit: 2, Quantum: 3, InitQuantum: 4, RateEnable: 5, FlowDefaultRate: 6, FlowMaxRate: 7, BucketsLog: 8, FlowRefillDelay: 9, OrphanMask: 10, LowRateThreshold: 11, CEThreshold: 12}}},
		"fq_codel": {val: &Attribute{Kind: "fq_codel", FqCodel: &FqCodel{Target: 1, Limit: 2, Interval: 3, ECN: 4, Flows: 5, Quantum: 6, CEThreshold: 7, DropBatchSize: 8, MemoryLimit: 9}}},
		"hfsc":     {val: &Attribute{Kind: "hfsc", Hfsc: &Hfsc{Rsc: &ServiceCurve{M1: 12, D: 34, M2: 56}}}},
		"hhf":      {val: &Attribute{Kind: "hhf", Hhf: &Hhf{BacklogLimit: 1, Quantum: 2, HHFlowsLimit: 3, ResetTimeout: 4, AdmitBytes: 5, EVICTTimeout: 6, NonHHWeight: 7}}},
		"htb":      {val: &Attribute{Kind: "htb", Htb: &Htb{Rate64: 123, Parms: &HtbOpt{Buffer: 0xFFFF}}}},
		"mqprio":   {val: &Attribute{Kind: "mqprio", MqPrio: &MqPrio{Mode: 1, Shaper: 2, MinRate64: 3, MaxRate64: 4}}},
		"pie":      {val: &Attribute{Kind: "pie", Pie: &Pie{Target: 1, Limit: 2, TUpdate: 3, Alpha: 4, Beta: 5, ECN: 6, Bytemode: 7}}},
		"qfq":      {val: &Attribute{Kind: "qfq", Qfq: &Qfq{Weight: 2, Lmax: 4}}},
		"red":      {val: &Attribute{Kind: "red", Red: &Red{MaxP: 2, Parms: &RedQOpt{QthMin: 2, QthMax: 4}}}},
		"sfb":      {val: &Attribute{Kind: "sfb", Sfb: &Sfb{Parms: &SfbQopt{Max: 0xFF}}}},
		"tbf":      {val: &Attribute{Kind: "tbf", Tbf: &Tbf{Rate64: 1, Prate64: 2, Burst: 3, Pburst: 4}}},
		"pfifo":    {val: &Attribute{Kind: "pfifo", Pfifo: &FifoOpt{Limit: 42}}},
		"bfifo":    {val: &Attribute{Kind: "bfifo", Bfifo: &FifoOpt{Limit: 84}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			options, err1 := validateQdiscObject(unix.RTM_NEWQDISC, &Object{Msg{Ifindex: 42}, *testcase.val})
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
