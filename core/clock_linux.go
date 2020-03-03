// +build linux

package core

import (
	"fmt"
	"os"
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
