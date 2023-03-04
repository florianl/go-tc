package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHfsc(t *testing.T) {
	tests := map[string]struct {
		val  Hfsc
		err1 error
		err2 error
	}{
		"Rsc": {val: Hfsc{Rsc: &ServiceCurve{M1: 12, D: 34, M2: 56}}},
		"Fsc": {val: Hfsc{Fsc: &ServiceCurve{M1: 13, D: 35, M2: 57}}},
		"Usc": {val: Hfsc{Usc: &ServiceCurve{M1: 14, D: 36, M2: 58}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHfsc(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Hfsc{}
			err2 := unmarshalHfsc(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Hfsc missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalHfsc(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestHfscQOpt(t *testing.T) {
	tests := map[string]struct {
		val  HfscQOpt
		err1 error
		err2 error
	}{
		"defcls 1":      {val: HfscQOpt{DefCls: 1}},
		"defcls 0":      {val: HfscQOpt{DefCls: 0}},
		"defcls 0xffff": {val: HfscQOpt{DefCls: 0xffff}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHfscQOpt(&testcase.val)
			if err1 != nil {
				t.Fatalf("unexpected error: %v", err1)
			}
			val := HfscQOpt{}
			err2 := unmarshalHfscQOpt(data, &val)
			if err2 != nil {
				t.Fatalf("unexpected error: %v", err2)
			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("HfscQOpt mismiatch (want+got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalHfscQOpt(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
