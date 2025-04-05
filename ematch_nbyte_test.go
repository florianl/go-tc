package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNByteMatch(t *testing.T) {
	tests := map[string]struct {
		val  NByteMatch
		err1 error
		err2 error
	}{
		"simple": {
			val: NByteMatch{Needle: []byte("helloWorld"),
				Offset: 42,
				Layer:  7},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalNByteMatch(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := NByteMatch{}
			err2 := unmarshalNByteMatch(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("NByteMatch missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalNByteMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUnmarshalNByteMatch(t *testing.T) {
	tests := map[string]struct {
		data      []byte
		needleLen uint16
		err       error
	}{
		"invalid length": {
			data: []byte{0x0, 0x1, 0x2, 0x3},
			err:  ErrInvalidArg,
		},
		"invalid needle": {
			data: []byte{0x00, 0x00,
				0xaa, 0xaa,
				0x01,
				0x00, 0x00, 0x00,
				0x0a, 0x0b, 0x0c},
			err: ErrInvalidArg,
		},
	}

	for name, testcase := range tests {
		name := name
		testcase := testcase
		t.Run(name, func(t *testing.T) {
			info := NByteMatch{}
			if err := unmarshalNByteMatch(testcase.data, &info); err != nil {
				if !errors.Is(err, ErrInvalidArg) {
					t.Fatalf("Unexpected error: %v", err)
				}
			}
		})
	}
}
