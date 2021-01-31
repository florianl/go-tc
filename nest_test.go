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
		"uint8":    {interpretation: vtUint8, attributeType: 1, data: uint8(123), result: []byte{0x5, 0x0, 0x1, 0x0, 0x7b, 0x0, 0x0, 0x0}},
		"uint16":   {interpretation: vtUint16, attributeType: 2, data: uint16(124), result: []byte{0x6, 0x0, 0x2, 0x0, 0x7c, 0x0, 0x0, 0x0}},
		"uint32":   {interpretation: vtUint32, attributeType: 3, data: uint32(125), result: []byte{0x8, 0x0, 0x3, 0x0, 0x7d, 0x0, 0x0, 0x0}},
		"uint64":   {interpretation: vtUint64, attributeType: 4, data: uint64(126), result: []byte{0xc, 0x0, 0x4, 0x0, 0x7e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		"string":   {interpretation: vtString, attributeType: 5, data: string("hello world"), result: []byte{0x10, 0x0, 0x5, 0x0, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0}},
		"bytes":    {interpretation: vtBytes, attributeType: 6, data: []byte{0x60, 0x0D, 0xCA, 0xFE}, result: []byte{0x8, 0x0, 0x6, 0x0, 0x60, 0xd, 0xca, 0xfe}},
		"flags":    {interpretation: vtFlag, attributeType: 7, data: []byte{}, result: []byte{0x04, 0x00, 0x07, 0x00}},
		"int8":     {interpretation: vtInt8, attributeType: 8, data: int8(-8), result: []byte{0x5, 0x0, 0x8, 0x0, 0xF8, 0x0, 0x0, 0x0}},
		"int16":    {interpretation: vtInt16, attributeType: 9, data: int16(-9), result: []byte{0x6, 0x0, 0x9, 0x0, 0xF7, 0xFF, 0x0, 0x0}},
		"int32":    {interpretation: vtInt32, attributeType: 10, data: int32(-10), result: []byte{0x8, 0x0, 0xA, 0x0, 0xF6, 0xFF, 0xFF, 0xFF}},
		"int64":    {interpretation: vtInt64, attributeType: 11, data: int64(-11), result: []byte{0xC, 0x0, 0xB, 0x0, 0xF5, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		"uint16Be": {interpretation: vtUint16Be, attributeType: 12, data: uint16(124), result: []byte{0x6, 0x0, 0xC, 0x0, 0x0, 0x7c, 0x0, 0x0}},
		"uint32Be": {interpretation: vtUint32Be, attributeType: 13, data: uint32(125), result: []byte{0x8, 0x0, 0xD, 0x0, 0x0, 0x0, 0x0, 0x7d}},
		"uint64Be": {interpretation: vtUint64Be, attributeType: 14, data: uint64(126), result: []byte{0xc, 0x0, 0xE, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7e}},
		"unknown":  {interpretation: vtUint64Be + 1, attributeType: 42, data: nil, err: fmt.Errorf("Unknown interpretation: 13")},
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
					t.Logf("recv: %v\n", err)
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

func TestUnmarshalAttributes(t *testing.T) {
	var valInt8 int8
	if err := unmarshalNetlinkAttribute([]byte{0xF8}, &valInt8); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valInt8 != -8 {
		t.Fatalf("expexted: -8\tgot: %d", valInt8)
	}

	var valInt16 int16
	if err := unmarshalNetlinkAttribute([]byte{0xF7, 0xFF}, &valInt16); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valInt16 != -9 {
		t.Fatalf("expexted: -8\tgot: %d", valInt8)
	}
	var valInt32 int32
	if err := unmarshalNetlinkAttribute([]byte{0xF6, 0xFF, 0xFF, 0xFF}, &valInt32); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valInt32 != -10 {
		t.Fatalf("expexted: -8\tgot: %d", valInt8)
	}
	var valInt64 int64
	if err := unmarshalNetlinkAttribute([]byte{0xF5, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, &valInt64); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valInt64 != -11 {
		t.Fatalf("expexted: -8\tgot: %d", valInt8)
	}

}
