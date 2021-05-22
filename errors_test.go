package tc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestConcatError(t *testing.T) {
	t.Run("nil + nil", func(t *testing.T) {
		result := concatError(nil, nil)
		if result != nil {
			t.Fatalf("expected nil but got %v", result)
		}
	})
	t.Run("io.EOF + nil", func(t *testing.T) {
		result := concatError(io.EOF, nil)
		if !errors.Is(result, io.EOF) {
			t.Fatalf("expected io.EOF but got %v", result)
		}
	})
	t.Run("nil + io.EOF", func(t *testing.T) {
		result := concatError(nil, io.EOF)
		if !errors.Is(result, io.EOF) {
			t.Fatalf("expected io.EOF but got %v", result)
		}
	})
	t.Run("io.EOF + os.ErrPermission", func(t *testing.T) {
		result := concatError(io.EOF, os.ErrPermission)
		fmt.Println(result)
		// Output:
		// EOF
		// permission denied
	})
}
