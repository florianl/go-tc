package core

import (
	"sync"
	"syscall"
	"time"
)

var (
	defaultClock Clock
	defaultOnce  sync.Once
)

const (
	// iproute2/include/utils.h:timeUnitsPerSec
	timeUnitsPerSec = 1000000
)

// Clock defines the interface for timing conversion operations
type Clock interface {
	// Duration2TcTime converts a given duration into a time value that can be converted
	// to ticks with Time2Tick().
	Duration2TcTime(d time.Duration) (uint32, error)

	// Time2Tick returns the number of CPU ticks for a given time in usec.
	Time2Tick(time uint32) uint32

	// Tick2Time returns a time in usec for a given number of CPU ticks.
	Tick2Time(tick uint32) uint32

	// XmitTime returns the time, that is needed to transmit a given size for a given rate.
	XmitTime(rate uint64, size uint32) uint32

	// XmitSize returns the size that can be transmitted at a given rate during a given time.
	XmitSize(rate uint64, ticks uint32) uint32

	// Time2Ktime converts time to kernel time units.
	Time2Ktime(time uint32) uint32

	// Ktime2Time converts kernel time units to time.
	Ktime2Time(ktime uint32) uint32

	// ClockFactor returns the current clock factor.
	ClockFactor() float64

	// TickInUSec returns the current tick in microseconds factor.
	TickInUSec() float64
}

// SystemClock implements the Clock interface using system-specific timing values
type SystemClock struct {
	clockFactor float64
	tickInUSec  float64
}

// NewSystemClock creates a new SystemClock instance by reading system timing values
func NewSystemClock() (*SystemClock, error) {
	cf, tick, err := readSystemClock()
	if err != nil {
		return nil, err
	}
	return &SystemClock{
		clockFactor: cf,
		tickInUSec:  tick,
	}, nil
}

// NewSystemClockWithDefaults creates a new SystemClock instance with default values
func NewSystemClockWithDefaults() *SystemClock {
	return &SystemClock{
		clockFactor: 1.0,
		tickInUSec:  1.0,
	}
}

// Duration2TcTime implements iproute2/tc/q_netem.c:get_ticks().
// It converts a given duration into a time value that can be converted to ticks with Time2Tick().
// On error it returns syscall.EINVAL.
func (c *SystemClock) Duration2TcTime(d time.Duration) (uint32, error) {
	v := uint64(int64(d.Microseconds()) * (timeUnitsPerSec / 1000000))
	if (v >> 32) != 0 {
		return 0, syscall.EINVAL
	}
	return uint32(v), nil
}

// Time2Tick implements iproute2/tc/tc_core:tc_core_time2tick().
// It returns the number of CPU ticks for a given time in usec.
func (c *SystemClock) Time2Tick(time uint32) uint32 {
	return uint32(float64(time) * c.tickInUSec)
}

// Tick2Time implements iproute2/tc/tc_core:tc_core_tick2time().
// It returns a time in usec for a given number of CPU ticks.
func (c *SystemClock) Tick2Time(tick uint32) uint32 {
	return uint32(float64(tick) / c.tickInUSec)
}

// XmitTime implements iproute2/tc/tc_core:tc_calc_xmittime().
// It returns the time, that is needed to transmit a given size for a given rate.
func (c *SystemClock) XmitTime(rate uint64, size uint32) uint32 {
	return c.Time2Tick(uint32(timeUnitsPerSec * (float64(size) / float64(rate))))
}

// XmitSize implements iproute2/tc/tc_core:tc_calc_xmitsize().
// It returns the size that can be transmitted at a given rate during a given time.
func (c *SystemClock) XmitSize(rate uint64, ticks uint32) uint32 {
	return uint32(rate*uint64(c.Tick2Time(ticks))) / timeUnitsPerSec
}

// Time2Ktime implements iproute2/tc/tc_core:tc_core_time2ktime().
func (c *SystemClock) Time2Ktime(time uint32) uint32 {
	return uint32(uint64(time) * uint64(c.clockFactor))
}

// Ktime2Time implements iproute2/tc/tc_core:tc_core_ktime2time().
func (c *SystemClock) Ktime2Time(ktime uint32) uint32 {
	return uint32(float64(ktime) / c.clockFactor)
}

// ClockFactor returns the current clock factor.
func (c *SystemClock) ClockFactor() float64 {
	return c.clockFactor
}

// TickInUSec returns the current tick in microseconds factor.
func (c *SystemClock) TickInUSec() float64 {
	return c.tickInUSec
}

// getDefaultClock returns the default clock instance, initializing it if necessary
func getDefaultClock() Clock {
	defaultOnce.Do(func() {
		clock, err := NewSystemClock()
		if err != nil {
			// Fall back to default values if we can't read system values
			defaultClock = NewSystemClockWithDefaults()
		} else {
			defaultClock = clock
		}
	})
	return defaultClock
}

// Duration2TcTime implements iproute2/tc/q_netem.c:get_ticks().
// It converts a given duration into a time value that can be converted to ticks with Time2Tick().
// On error it returns syscall.EINVAL.
// Deprecated: Use (*SystemClock).Duration2TcTime instead.
func Duration2TcTime(d time.Duration) (uint32, error) {
	return getDefaultClock().Duration2TcTime(d)
}

// Time2Tick implements iproute2/tc/tc_core:tc_core_time2tick().
// It returns the number of CPU ticks for a given time in usec.
// Deprecated: Use (*SystemClock).Time2Tick instead.
func Time2Tick(time uint32) uint32 {
	return getDefaultClock().Time2Tick(time)
}

// Tick2Time implements iproute2/tc/tc_core:tc_core_tick2time().
// It returns a time in usec for a given number of CPU ticks.
// Deprecated: Use (*SystemClock).Tick2Time instead.
func Tick2Time(tick uint32) uint32 {
	return getDefaultClock().Tick2Time(tick)
}

// XmitTime implements iproute2/tc/tc_core:tc_calc_xmittime().
// It returns the time, that is needed to transmit a given size for a given rate.
// Deprecated: Use (*SystemClock).XmitTime instead.
func XmitTime(rate uint64, size uint32) uint32 {
	return getDefaultClock().XmitTime(rate, size)
}

// XmitSize implements iproute2/tc/tc_core:tc_calc_xmitsize().
// It returns the size that can be transmitted at a given rate during a given time.
// Deprecated: Use (*SystemClock).XmitSize instead.
func XmitSize(rate uint64, ticks uint32) uint32 {
	return getDefaultClock().XmitSize(rate, ticks)
}

// Time2Ktime implements iproute2/tc/tc_core:tc_core_time2ktime().
// Deprecated: Use (*SystemClock).Time2Ktime instead.
func Time2Ktime(time uint32) uint32 {
	return getDefaultClock().Time2Ktime(time)
}

// Ktime2Time implements iproute2/tc/tc_core:tc_core_ktime2time().
// Deprecated: Use (*SystemClock).Ktime2Time instead.
func Ktime2Time(ktime uint32) uint32 {
	return getDefaultClock().Ktime2Time(ktime)
}

// ClockFactor returns the current clock factor from the default clock instance.
// Deprecated: Use (*SystemClock).ClockFactor instead.
func ClockFactor() float64 {
	return getDefaultClock().ClockFactor()
}

// TickInUSec returns the current tick in microseconds factor from the default clock instance.
// Deprecated: Use (*SystemClock).TickInUSec instead.
func TickInUSec() float64 {
	return getDefaultClock().TickInUSec()
}
