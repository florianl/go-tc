package tc

import (
	"errors"
	"testing"
)

func TestPlug(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		_, err := marshalPlug(&Plug{
			Action: PlugLimit,
			Limit:  123,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("nil", func(t *testing.T) {
		_, err := marshalPlug(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unmarshal", func(t *testing.T) {
		err := unmarshalPlug([]byte{}, &Plug{})
		if !errors.Is(err, ErrNotImplemented) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
