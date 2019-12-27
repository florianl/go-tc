package core

import (
	"fmt"
	"os"
)

var tickInUSec uint32
var clockFactor uint32

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
	fmt.Println(t2us, us2t, clockRes, hiClockRes)

	clockFactor = uint32(float64(clockRes) / float64(timeUnitsPerSec))
	tickInUSec = uint32(float64(t2us/us2t) * float64(clockFactor))
}

// iproute2/tc/tc_core:tc_core_time2tick()
func time2tick(time uint32) uint32 {
	return time * tickInUSec
}

// iproute2/tc/tc_core:tc_core_tick2time()
func tick2time(tick uint32) uint32 {
	return tick / tickInUSec
}

// iproute2/tc/tc_core:tc_calc_xmittime()
func xmittime(rate uint64, size uint32) uint32 {
	return time2tick(uint32(timeUnitsPerSec * (float64(size) / float64(rate))))

}

// iproute2/tc/tc_core:tc_calc_xmitsize()
func xmitsize(rate uint64, ticks uint32) uint32 {
	return uint32(rate*uint64(tick2time(ticks))) / timeUnitsPerSec
}
