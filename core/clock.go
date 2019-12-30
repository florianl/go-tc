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

// iproute2/tc/tc_core:tc_core_time2tick()
func time2tick(time uint32) uint32 {
	return uint32(float64(time) * tickInUSec)
}

// iproute2/tc/tc_core:tc_core_tick2time()
func tick2time(tick uint32) uint32 {
	return tick / uint32(tickInUSec)
}

// iproute2/tc/tc_core:tc_calc_xmittime()
func xmittime(rate uint64, size uint32) uint32 {
	return time2tick(uint32(timeUnitsPerSec * (float64(size) / float64(rate))))

}

// iproute2/tc/tc_core:tc_calc_xmitsize()
func xmitsize(rate uint64, ticks uint32) uint32 {
	return uint32(rate*uint64(tick2time(ticks))) / timeUnitsPerSec
}
