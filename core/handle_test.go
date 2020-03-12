package core

import "testing"

// Tests out the HandleStr function
func TestSplitHandle(t *testing.T) {
	tests := []struct {
		name  string
		args  uint32
		major uint32
		minor uint32
	}{
		{"handle 0", 0, 0, 0},
		{"handle 65535", 0x0000ffff, 0, 65535},
		{"handle 4294901760", 0xffff0000, 65535, 0},
		{"handle 4294967295", 0xffffffff, 65535, 65535},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if maj, min := SplitHandle(tt.args); maj != tt.major && min != tt.minor {
				t.Errorf("HandleStr() = %d:%d, want %d:%d", maj, min, tt.major, tt.minor)
			}
		})
	}
}

// Test the BuildHandleFunction
func TestBuildHandle(t *testing.T) {
	type args struct {
		maj uint32
		min uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{"0:2", args{0, 2}, 2},
		{"0:65535", args{0, 65535}, 0x0000ffff},
		{"65535:65535", args{65535, 65535}, 0xffffffff},
		{"65535:65535", args{65535, 0}, 0xffff0000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildHandle(tt.args.maj, tt.args.min); got != tt.want {
				t.Errorf("BuildHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
