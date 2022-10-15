package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIptMatch(t *testing.T) {
	tests := map[string]struct {
		val  IptMatch
		err1 error
		err2 error
	}{
		"simple": {val: IptMatch{
			Hook:      uint32Ptr(1),
			MatchName: stringPtr("foo"),
			Revision:  uint8Ptr(2),
			NFProto:   uint8Ptr(3),
			MatchData: bytesPtr([]byte{0xAA, 0x55, 0xAA, 0x55}),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalIptMatch(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := IptMatch{}
			err2 := unmarshalIptMatch(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("IptMatch missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalIptMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
