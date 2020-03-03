package core

var tickInUSec float64
var clockFactor float64

const (
	// iproute2/include/utils.h:timeUnitsPerSec
	timeUnitsPerSec = 1000000
)

// Time2Tick implements iproute2/tc/tc_core:tc_core_time2tick().
// It returns the number of CPU ticks for a given time in usec.
func Time2Tick(time uint32) uint32 {
	return uint32(float64(time) * tickInUSec)
}

// Tick2Time implements iproute2/tc/tc_core:tc_core_tick2time().
// It returns a time in usec for a given number of CPU ticks.
func Tick2Time(tick uint32) uint32 {
	return uint32(float64(tick) / tickInUSec)
}

// XmitTime implements iproute2/tc/tc_core:tc_calc_xmittime().
// It returns the time, that is needed to transmit a given size for a given rate.
func XmitTime(rate uint64, size uint32) uint32 {
	return Time2Tick(uint32(timeUnitsPerSec * (float64(size) / float64(rate))))

}

// XmitSize implements iproute2/tc/tc_core:tc_calc_xmitsize().
// It returns the size that can be transmitted at a given rate during a given time.
func XmitSize(rate uint64, ticks uint32) uint32 {
	return uint32(rate*uint64(Tick2Time(ticks))) / timeUnitsPerSec
}

// Time2Ktime implements iproute2/tc/tc_core:tc_core_time2ktime().
func Time2Ktime(time uint32) uint32 {
	return uint32(uint64(time) * uint64(clockFactor))
}

// Ktime2Time implements iproute2/tc/tc_core:tc_core_ktime2time().
func Ktime2Time(ktime uint32) uint32 {
	return uint32(float64(ktime) / clockFactor)
}
