package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFqPie(t *testing.T) {
	tests := map[string]struct {
		val  FqPie
		err1 error
		err2 error
	}{
		"simple": {val: FqPie{ // Used defaults from https://man7.org/linux/man-pages/man8/tc-fq_pie.8.html
			Limit:           uint32Ptr(10240),
			Flows:           uint32Ptr(1024),
			Target:          uint32Ptr(15),
			TUpdate:         uint32Ptr(15),
			Alpha:           uint32Ptr(2),
			Beta:            uint32Ptr(20),
			Quantum:         uint32Ptr(1514),
			MemoryLimit:     uint32Ptr(32),
			EcnProb:         uint32Ptr(10),
			Ecn:             uint32Ptr(0),
			Bytemode:        uint32Ptr(0),
			DqRateEstimator: uint32Ptr(0)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFqPie(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := FqPie{}
			err2 := unmarshalFqPie(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("FqPie mismatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalFqPie(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("nil", func(t *testing.T) {
		err := unmarshalFqPie([]byte{0x0}, nil)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
	})
}
