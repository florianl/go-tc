//+build !linux

package tc

// Open establishes a RTNETLINK socket for traffic control
func Open(config *Config) (*Tc, error) { return &Tc{}, ErrNotLinux }

func (tc *Tc) Close() error { return ErrNotLinux }
