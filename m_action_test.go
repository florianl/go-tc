package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAction(t *testing.T) {
	tests := map[string]struct {
		val  Action
		err1 error
		err2 error
	}{
		"empty":               {err1: fmt.Errorf("kind is missing")},
		"unknown Kind":        {val: Action{Kind: "test", Index: 123}, err1: fmt.Errorf("unknown kind 'test'")},
		"bpf Without Options": {val: Action{Kind: "bpf"}, err1: fmt.Errorf("ActBpf options are missing")},
		"simple Bpf":          {val: Action{Kind: "bpf", Bpf: &ActBpf{FD: 12, Name: "simpleTest", Parms: &ActBpfParms{Action: 2, Index: 4}}}},
		"connmark":            {val: Action{Kind: "connmark", ConnMark: &Connmark{Parms: &ConnmarkParam{Index: 42, Action: 1}}}},
		"csum":                {val: Action{Kind: "csum", CSum: &Csum{Parms: &CsumParms{Index: 1, Capab: 2}}}},
		"defact":              {val: Action{Kind: "defact", Defact: &Defact{Parms: &DefactParms{Index: 42, Action: 1}}}},
		"ife":                 {val: Action{Kind: "ife", Ife: &Ife{Parms: &IfeParms{Index: 42, Action: 1}}}},
		"ipt":                 {val: Action{Kind: "ipt", Ipt: &Ipt{Table: "testTable", Hook: 42, Index: 1984}}},
		"mirred":              {val: Action{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 42, Action: 1}}}},
		"nat":                 {val: Action{Kind: "nat", Nat: &Nat{Parms: &NatParms{Index: 42, Action: 1}}}},
		"police":              {val: Action{Kind: "police", Police: &Police{AvRate: 1337, Result: 42}}},
		"sample":              {val: Action{Kind: "sample", Sample: &Sample{Parms: &SampleParms{Index: 42, Action: 1}}}},
		"vlan":                {val: Action{Kind: "vlan", VLan: &VLan{Parms: &VLanParms{Index: 42, Action: 1}}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalAction(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Action{}
			err2 := unmarshalAction(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Action missmatch (want +got):\n%s", diff)
			}
		})
	}
}
