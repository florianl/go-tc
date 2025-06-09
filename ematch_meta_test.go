package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMetaMatch(t *testing.T) {
	tests := map[string]struct {
		val  MetaMatch
		err1 error
		err2 error
	}{
		"simple": {val: MetaMatch{
			Hdr: &MetaHdr{
				Left: MetaValue{
					Kind:  1<<12 | 1,
					Shift: 2,
					Op:    3,
				},
				Right: MetaValue{
					Kind:  1<<12 | 4,
					Shift: 5,
					Op:    6,
				},
			},
			Left: &MetaValueType{
				Int: uint32Ptr(42),
			},
			Right: &MetaValueType{
				Int: uint32Ptr(73),
			},
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMetaMatch(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := MetaMatch{}
			err2 := unmarshalMetaMatch(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("MetaMatch missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalMetaMatch(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unmarshalIptMatch()", func(t *testing.T) {
		err := unmarshalMetaMatch([]byte{0x0}, nil)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}
