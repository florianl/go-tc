package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRsvp(t *testing.T) {
	actions := []*Action{
		{Kind: "csum", CSum: &Csum{Parms: &CsumParms{Index: 4, Capab: 5}}},
	}

	tests := map[string]struct {
		val  Rsvp
		err1 error
		err2 error
	}{
		"simple":      {val: Rsvp{ClassID: uint32Ptr(43), Src: bytesPtr([]byte{0xAA}), Dst: bytesPtr([]byte{0x55}), Police: &Police{AvRate: uint32Ptr(1337), Result: uint32Ptr(12)}}},
		"with Action": {val: Rsvp{ClassID: uint32Ptr(73), Actions: &actions}},
		"extended":    {val: Rsvp{ClassID: uint32Ptr(13), Src: bytesPtr([]byte{0xAA}), Dst: bytesPtr([]byte{0x55}), PInfo: &RsvpPInfo{Dpi: RsvpGpi{Mask: 1234, Key: 4321, Offset: 1}, Protocol: 42}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalRsvp(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Rsvp{}
			err2 := unmarshalRsvp(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Rsvp missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalRsvp(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
