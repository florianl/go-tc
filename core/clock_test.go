package core

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	t.Run("Tick", func(t *testing.T) {
		tick := Time2Tick(0xC0FFEE)
		time := Tick2Time(tick)
		// to avoid the dependency to math and make sure, diff is always equal or
		// bigger than 0, diff is the square of the difference of both values
		diff := (0xC0FFEE - time) * (0xC0FFEE - time)
		if diff > 3 {
			t.Fatalf("expected %d, got %d", 0xC0FFEE, time)
		}
	})

	t.Run("Ktime", func(t *testing.T) {
		ktime := Time2Ktime(0xC0FFEE)
		time := Ktime2Time(ktime)
		diff := (0xC0FFEE - time) * (0xC0FFEE - time)
		if diff > 3 {
			t.Fatalf("expected %d, got %d", 0xC0FFEE, time)
		}
	})

	t.Run("XmitTime", func(t *testing.T) {
		timeA := XmitTime(4096, 4096)
		timeB := Time2Tick(timeUnitsPerSec)
		if timeA != timeB {
			t.Fatalf("expected %d, got %d", timeB, timeA)
		}
	})
	t.Run("XmitSize", func(t *testing.T) {
		a := XmitSize(timeUnitsPerSec, 1)
		b := Tick2Time(1)
		if a != b {
			t.Fatalf("expected %d, got %d", a, b)
		}
	})
}

func TestDuration2TcTime(t *testing.T) {
	tests := map[string]struct {
		d    time.Duration
		time uint32
		err  error
	}{
		"73 m":  {d: 73 * time.Minute, err: syscall.EINVAL},
		"73 s":  {d: 73 * time.Second, time: 73000000},
		"73 ms": {d: 73 * time.Millisecond, time: 73000},
		"73 us": {d: 73 * time.Microsecond, time: 73},
	}

	for name, testcase := range tests {
		name := name
		testcase := testcase
		t.Run(name, func(t *testing.T) {
			time, err := Duration2TcTime(testcase.d)
			if !errors.Is(err, testcase.err) {
				t.Fatalf("Expected %v but got %v", testcase.err, err)
			}
			if time != testcase.time {
				t.Fatalf("Expected %d but got %d", testcase.time, time)
			}
		})
	}
}
