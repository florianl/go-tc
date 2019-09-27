package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFw(t *testing.T) {
	tests := map[string]struct {
		val  Fw
		err1 error
		err2 error
	}{
		"empty":    {},
		"simple":   {val: Fw{ClassID: 12, InDev: "lo", Mask: 0xFFFF}},
		"extended": {val: Fw{ClassID: 12, InDev: "lo", Mask: 0xFFFF, Police: &Police{AvRate: 1337, Result: 12}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFw(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Fw{}
			err2 := unmarshalFw(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Fw missmatch (want +got):\n%s", diff)
			}
		})
	}
}
