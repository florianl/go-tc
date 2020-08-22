package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRsvp(t *testing.T) {
	tests := map[string]struct {
		val  Rsvp
		err1 error
		err2 error
	}{
		"simple":   {val: Rsvp{ClassID: uint32Ptr(43), Src: bytesPtr([]byte{0xAA}), Dst: bytesPtr([]byte{0x55}), Police: &Police{AvRate: 1337, Result: 12}}},
		"extended": {val: Rsvp{ClassID: uint32Ptr(13), Src: bytesPtr([]byte{0xAA}), Dst: bytesPtr([]byte{0x55}), PInfo: &RsvpPInfo{Dpi: RsvpGpi{Mask: 1234, Key: 4321, Offset: 1}, Protocol: 42}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalRsvp(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Rsvp{}
			err2 := unmarshalRsvp(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
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
