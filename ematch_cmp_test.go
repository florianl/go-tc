package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmpMatch(t *testing.T) {
	// cmp(u16 at 3 layer 2 mask 0xff00 gt 20)
	in := CmpMatch{
		Val:   20,
		Mask:  0xff00,
		Off:   3,
		Align: CmpMatchU16,
		Layer: EmatchLayerTransport,
		Opnd:  EmatchOpndGt,
	}

	data, err := marshalCmpMatch(&in)
	if err != nil {
		t.Fatal(err)
	}
	out := CmpMatch{}
	if err := unmarshalCmpMatch(data, &out); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(in, out); diff != "" {
		t.Fatalf("CmpMatch missmatch (-want +got):\n%s", diff)
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalCmpMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
