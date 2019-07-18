package tc

import (
	"testing"

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
			if err := extractTcmsgAttributes(testcase.input, value); err != nil {
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
