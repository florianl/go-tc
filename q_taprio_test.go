package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTaPrio(t *testing.T) {
	tests := map[string]struct {
		val  TaPrio
		err1 error
		err2 error
	}{
		"simple": {val: TaPrio{
			PrioMap: &MqPrioQopt{
				NumTc: 3},
			SchedBaseTime:           int64Ptr(5),
			SchedClockID:            int32Ptr(7),
			SchedCycleTime:          int64Ptr(11),
			SchedCycleTimeExtension: int64Ptr(13),
			Flags:                   uint32Ptr(17),
			TxTimeDelay:             uint32Ptr(19),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalTaPrio(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := TaPrio{}
			err2 := unmarshalTaPrio(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Etf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalTaPrio(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
