package tc

import (
	"math"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertNetIP(t *testing.T) {
	tests := map[string]struct {
		value net.IP
	}{
		"ipv6-localhost": {value: net.ParseIP("::1")},
		"ipv4-localhost": {value: net.ParseIP("127.0.0.1")},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := netIPPtr(testcase.value)
			value := netIPValue(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertNetIPMask(t *testing.T) {
	tests := map[string]struct {
		value net.IPMask
	}{
		"broadcast":    {value: net.IPv4Mask(255, 255, 255, 255)},
		"default mask": {value: net.ParseIP("127.0.0.1").DefaultMask()},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := netIPMaskPtr(testcase.value)
			value := netIPMaskValue(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertString(t *testing.T) {
	tests := map[string]struct {
		value string
	}{
		"hello world":  {value: "hello world"},
		"empty string": {value: ""},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := stringPtr(testcase.value)
			value := stringValue(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertUint8(t *testing.T) {
	tests := map[string]struct {
		value uint8
	}{
		"0":         {value: 0},
		"uint8 max": {value: math.MaxUint8},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := uint8Ptr(testcase.value)
			value := uint8Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertUint16(t *testing.T) {
	tests := map[string]struct {
		value uint16
	}{
		"0":          {value: 0},
		"uint16 max": {value: math.MaxUint16},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := uint16Ptr(testcase.value)
			value := uint16Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertUint32(t *testing.T) {
	tests := map[string]struct {
		value uint32
	}{
		"0":          {value: 0},
		"uint32 max": {value: math.MaxUint32},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := uint32Ptr(testcase.value)
			value := uint32Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertUint64(t *testing.T) {
	tests := map[string]struct {
		value uint64
	}{
		"0":          {value: 0},
		"uint64 max": {value: math.MaxUint64},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := uint64Ptr(testcase.value)
			value := uint64Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInt8(t *testing.T) {
	tests := map[string]struct {
		value int8
	}{
		"0":        {value: 0},
		"int8 max": {value: math.MaxInt8},
		"int8 min": {value: math.MinInt8},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := int8Ptr(testcase.value)
			value := int8Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInt16(t *testing.T) {
	tests := map[string]struct {
		value int16
	}{
		"0":         {value: 0},
		"int16 max": {value: math.MaxInt16},
		"int16 min": {value: math.MinInt16},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := int16Ptr(testcase.value)
			value := int16Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInt32(t *testing.T) {
	tests := map[string]struct {
		value int32
	}{
		"0":         {value: 0},
		"int32 max": {value: math.MaxInt32},
		"int32 min": {value: math.MinInt32},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := int32Ptr(testcase.value)
			value := int32Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInt64(t *testing.T) {
	tests := map[string]struct {
		value int64
	}{
		"0":         {value: 0},
		"int64 max": {value: math.MaxInt64},
		"int64 min": {value: math.MinInt64},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := int64Ptr(testcase.value)
			value := int64Value(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertBytes(t *testing.T) {
	tests := map[string]struct {
		value []byte
	}{
		"empty":       {value: []byte{}},
		"single byte": {value: []byte{0xFF}},
		"disk sync":   {value: []byte{0xAA, 0x55}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := bytesPtr(testcase.value)
			value := bytesValue(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertBool(t *testing.T) {
	tests := map[string]struct {
		value bool
	}{
		"true":  {value: true},
		"false": {value: false},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := boolPtr(testcase.value)
			value := boolValue(ptr)
			if diff := cmp.Diff(value, testcase.value); diff != "" {
				t.Fatalf("Missmatch (-want +got):\n%s", diff)
			}
		})
	}
}
