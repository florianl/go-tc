//go:build go1.18
// +build go1.18

package tc

import (
	"reflect"
	"testing"
)

func FuzzU32(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var info U32
		if err := unmarshalU32(data, &info); err != nil {
			t.Skip()
		}

		new, err := marshalU32(&info)
		if err != nil {
			t.Fatalf("failed to marshal %#v: %v", info, err)
		}
		if !reflect.DeepEqual(data, new) {
			t.Errorf("(un-)marshal missmatch:\n%v\n%v", data, new)
		}
	})
}
