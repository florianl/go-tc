//+build !linux

package tc

import "testing"

func TestOthers(t *testing.T) {
	socket := &Tc{}
	want := ErrNotLinux

	if _, got := Open(&Config{}); got != want {
		t.Fatalf("unexpected error:\ngot:\t%v\nwant:\t%v\n", got, want)
	}

	if got := socket.Close(); got != want {
		t.Fatalf("unexpected error:\ngot:\t%v\nwant:\t%v\n", got, want)
	}
}
