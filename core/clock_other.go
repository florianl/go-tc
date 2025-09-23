//go:build !linux
// +build !linux

package core

// readSystemClock returns default clock values for non-Linux platforms
func readSystemClock() (float64, float64, error) {
	return 1.0, 1.0, nil
}
