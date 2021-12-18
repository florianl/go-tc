package tc

import (
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestU32Match(t *testing.T) {
	in := U32Match{
		Mask:  0xff,
		Value: 0xaa,
	}

	data, err := marshalU32Match(&in)
	if err != nil {
		t.Fatal(err)
	}
	out := U32Match{}
	if err := unmarshalU32Match(data, &out); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(in, out); diff != "" {
		t.Fatalf("U32Match missmatch (-want +got):\n%s", diff)
	}
	t.Run("nil-marshalU32Match", func(t *testing.T) {
		_, err := marshalU32Match(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("nil-unmarshalU32Match", func(t *testing.T) {
		if err := unmarshalU32Match([]byte{}, &U32Match{}); !errors.Is(err, io.EOF) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
