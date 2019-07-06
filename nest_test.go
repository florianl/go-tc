package tc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMarshalAttributes(t *testing.T) {
	tests := map[string]struct {
		interpretation valueType
		attributeType  uint16
		data           interface{}
		result         []byte
		err            error
	}{
		"uint8":   {interpretation: vtUint8, attributeType: 1, data: uint8(123), result: []byte{0x5, 0x0, 0x1, 0x0, 0x7b, 0x0, 0x0, 0x0}},
		"uint16":  {interpretation: vtUint16, attributeType: 2, data: uint16(124), result: []byte{0x6, 0x0, 0x2, 0x0, 0x7c, 0x0, 0x0, 0x0}},
		"uint32":  {interpretation: vtUint32, attributeType: 3, data: uint32(125), result: []byte{0x8, 0x0, 0x3, 0x0, 0x7d, 0x0, 0x0, 0x0}},
		"uint64":  {interpretation: vtUint64, attributeType: 4, data: uint64(126), result: []byte{0xc, 0x0, 0x4, 0x0, 0x7e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		"string":  {interpretation: vtString, attributeType: 5, data: string("hello world"), result: []byte{0x10, 0x0, 0x5, 0x0, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0}},
		"bytes":   {interpretation: vtBytes, attributeType: 6, data: []byte{0x60, 0x0D, 0xCA, 0xFE}, result: []byte{0x8, 0x0, 0x6, 0x0, 0x60, 0xd, 0xca, 0xfe}},
		"unknown": {interpretation: vtBytes + 1, attributeType: 42, data: nil, err: fmt.Errorf("Unknown interpretation: 6")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := marshalAttributes([]tcOption{
				{
					Interpretation: tc.interpretation,
					Type:           tc.attributeType,
					Data:           tc.data,
				},
			})
			if err != nil {
				if tc.err != nil {
					// TODO compare resulting errors
					return
				}
				t.Fatalf("expexted no error and got: %v", err)
			}
			if tc.err != nil {
				t.Fatalf("expected error: %v", tc.err)
			}
			if !bytes.Equal(tc.result, got) {
				t.Fatalf("expected: %v\ngot: %v", tc.result, got)
			}
		})
	}
}
