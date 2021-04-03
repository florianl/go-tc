package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActBpft(t *testing.T) {
	tests := map[string]struct {
		val    ActBpf
		enrich *Tcft
		err1   error
		err2   error
	}{
		"simple":          {val: ActBpf{FD: uint32Ptr(12), Name: stringPtr("simpleTest")}},
		"invalidArgument": {val: ActBpf{FD: uint32Ptr(12), Name: stringPtr("simpleTest"), Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
		"extended":        {val: ActBpf{FD: uint32Ptr(12), Name: stringPtr("simpleTest"), Parms: &ActBpfParms{Action: 2, Index: 4}}},
		"Tm Attribute": {val: ActBpf{FD: uint32Ptr(12), Name: stringPtr("simpleTest"), Parms: &ActBpfParms{Action: 2, Index: 4}},
			enrich: &Tcft{Install: 1, LastUse: 2, Expires: 3, FirstUse: 4}},
		"legacy BPF": {val: ActBpf{Ops: bytesPtr([]byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}), OpsLen: uint16Ptr(1),
			ID: uint32Ptr(42), Tag: bytesPtr([]byte("foo"))}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalActBpf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			if testcase.enrich != nil {
				enrichment, err := marshalStruct(testcase.enrich)
				if err != nil {
					t.Fatalf("could not generate enrichment: %v", err)
				}
				tmp, _ := marshalAttributes([]tcOption{{
					Interpretation: vtBytes, Type: tcaActBpfTm, Data: enrichment}})
				data = append(data, tmp...)
				testcase.val.Tm = testcase.enrich
			}
			newData, tm := injectTcft(t, data, tcaActBpfTm)
			newData = injectAttribute(t, newData, []byte{}, tcaActBpfPad)
			val := ActBpf{}
			err2 := unmarshalActBpf(newData, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			testcase.val.Tm = tm
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("ActBpft missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalActBpf(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
