package tc

import (
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContainerMatch(t *testing.T) {
	tests := map[string]uint32{
		"TestPos5":   uint32(5),
		"TestPos100": uint32(100),
	}
	for name, pos := range tests {
		t.Run(name, func(t *testing.T) {
			in := ContainerMatch{
				Pos: pos,
			}

			data, err := marshalContainerMatch(&in)
			if err != nil {
				t.Fatal(err)
			}
			out := ContainerMatch{}
			if err := unmarshalContainerMatch(data, &out); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(in, out); diff != "" {
				t.Fatalf("ContainerMatch missmatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("invalid pos", func(t *testing.T) {
		in := ContainerMatch{
			Pos: 0,
		}
		if _, err := marshalContainerMatch(&in); !errors.Is(err, ErrInvalidArg) {
			t.Fatalf("Expected ErrInvalidArg but got '%v'", err)
		}
	})
	t.Run("nil-marshalContainerMatch", func(t *testing.T) {
		_, err := marshalContainerMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil-unmarshalContainerMatch", func(t *testing.T) {
		inByte := []byte{0x0001}

		if err := unmarshalContainerMatch(inByte, nil); !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil-byte-unmarshalContainerMatch", func(t *testing.T) {
		out := ContainerMatch{
			Pos: 10,
		}
		if err := unmarshalContainerMatch([]byte{}, &out); !errors.Is(err, io.EOF) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
