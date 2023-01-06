//go:build linux && amd64 && go1.17
// +build linux,amd64,go1.17

package core

import (
	"fmt"
	"os"
	"testing"
)

func TestReadPsched(t *testing.T) {
	t.Run("Read from system", func(t *testing.T) {
		if _, _, err := readPsched(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid psched format", func(t *testing.T) {
		tmpDir := t.TempDir()
		defer os.RemoveAll(tmpDir)

		if err := os.MkdirAll(fmt.Sprintf("%s/net", tmpDir), 0750); err != nil {
			t.Fatal(err)
		}
		fakePsched, err := os.Create(fmt.Sprintf("%s/net/psched", tmpDir))
		if err != nil {
			t.Fatal(err)
		}
		// Write some invalid data into the file
		if _, err := fakePsched.WriteString("hello world"); err != nil {
			t.Fatal(err)
		}

		t.Setenv("PROC_ROOT", tmpDir)

		if _, _, err := readPsched(); err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("psched does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		defer os.RemoveAll(tmpDir)
		t.Setenv("PROC_ROOT", tmpDir)

		if _, _, err := readPsched(); err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
