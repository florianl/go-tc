package tc

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMatchall(t *testing.T) {
	actions := []*Action{
		{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 0x1, Capab: 0x0, Action: 0x4, RefCnt: 0x1, BindCnt: 0x1, Eaction: 0x1, IfIndex: 0x2}}},
	}

	tests := map[string]struct {
		val  Matchall
		err1 error
		err2 error
	}{
		"simple":  {val: Matchall{ClassID: uint32Ptr(42), Flags: uint32Ptr(SkipHw)}},
		"actions": {val: Matchall{ClassID: uint32Ptr(1337), Flags: uint32Ptr(SkipHw), Actions: &actions}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalMatchall(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			pcnt := uint64(1337)
			pcntBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(pcntBytes, pcnt)
			tmp := injectAttribute(t, data, pcntBytes, tcaMatchallPcnt)
			newData := injectAttribute(t, tmp, []byte{}, tcaMatchallPad)
			val := Matchall{}
			err2 := unmarshalMatchall(newData, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}

			// Reinject pcnt
			testcase.val.Pcnt = uint64Ptr(pcnt)
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Matchall missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalMatchall(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unmarshalMatchall()", func(t *testing.T) {
		err := unmarshalMatchall([]byte{0x0}, nil)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}
