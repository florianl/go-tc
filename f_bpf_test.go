package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBpf(t *testing.T) {
	tests := map[string]struct {
		val  Bpf
		err1 error
		err2 error
	}{
		"empty": {err1: fmt.Errorf("Bpf options are missing")},
		"simple": {val: Bpf{Ops: []byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff},
			OpsLen:  0x1,
			ClassID: 0x10001,
			Flags:   0x1}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalBpf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Bpf{}
			err2 := unmarshalBpf(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Bpf missmatch (want +got):\n%s", diff)
			}
		})
	}
}
