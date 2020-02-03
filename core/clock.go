package core

import (
	"fmt"
	"os"
)

var tickInUSec float64
var clockFactor float64

const (
	// iproute2/include/utils.h:timeUnitsPerSec
	timeUnitsPerSec = 1000000
)

func init() {

	fd, err := os.Open("/proc/net/psched")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open /proc/net/psched: %v\n", err)
		return
	}
	defer fd.Close()

	var t2us, us2t, clockRes, hiClockRes uint32
	_, err = fmt.Fscanf(fd, "%08x %08x %08x %08x", &t2us, &us2t, &clockRes, &hiClockRes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read /proc/net/psched: %v\n", err)
	}
	clockFactor = float64(clockRes) / timeUnitsPerSec
	tickInUSec = float64(t2us) / float64(us2t) * clockFactor
}

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
