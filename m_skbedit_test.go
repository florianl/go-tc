package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSkbEdit(t *testing.T) {
	tests := map[string]struct {
		val  SkbEdit
		err1 error
		err2 error
	}{
		"simple": {val: SkbEdit{Parms: &SkbEditParms{BindCnt: 111}}},
		"all arguments": {val: SkbEdit{Parms: &SkbEditParms{Index: 222},
			Priority: uint32Ptr(11), QueueMapping: uint16Ptr(12), Mark: uint32Ptr(13), Ptype: uint16Ptr(14),
			Mask: uint32Ptr(15), Flags: uint64Ptr(16), QueueMappingMax: uint16Ptr(17)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalSkbEdit(&testcase.val)
			if !errors.Is(err1, testcase.err1) {
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := SkbEdit{}
			err2 := unmarshalSkbEdit(data, &val)
			if !errors.Is(err2, testcase.err2) {
				t.Fatalf("Unexpected error: %v", err2)
			}

			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("SkbEdit missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalSkbEdit(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
