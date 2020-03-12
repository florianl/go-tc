package core

import "testing"

// Tests out the HandleStr function
func TestSplitHandle(t *testing.T) {
	tests := map[string]struct {
		args  uint32
		major uint32
		minor uint32
	}{
		"handle 0":          {args: 0, major: 0, minor: 0},
		"handle 65535":      {args: 0x0000ffff, major: 0, minor: 65535},
		"handle 4294901760": {args: 0xffff0000, major: 65535, minor: 0},
		"handle 4294967295": {args: 0xffffffff, major: 65535, minor: 65535},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if maj, min := SplitHandle(tt.args); maj != tt.major && min != tt.minor {
				t.Errorf("HandleStr() = %d:%d, want %d:%d", maj, min, tt.major, tt.minor)
			}
		})
	}
}

// Test the BuildHandleFunction
func TestBuildHandle(t *testing.T) {
	tests := map[string]struct {
		major uint32
		minor uint32
		want  uint32
	}{
		"0:2":         {major: 0, minor: 2, want: 2},
		"0:65535":     {major: 0, minor: 65535, want: 0x0000ffff},
		"65535:65535": {major: 65535, minor: 65535, want: 0xffffffff},
		"65535:0":     {major: 65535, minor: 0, want: 0xffff0000},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := BuildHandle(tt.major, tt.minor); got != tt.want {
				t.Errorf("BuildHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
