package tc

import (
	"github.com/florianl/go-tc/internal/unix"
)

// Chain represents a collection of filter
type Chain struct {
	Tc
}

// Chain allows to read and alter chains
func (tc *Tc) Chain() *Chain {
	return &Chain{*tc}
}

func (c *Chain) Add(info *Object) error {
	return ErrNotImplemented
}

func (c *Chain) Delete() error {
	return ErrNotImplemented
}

func (c *Chain) Get() ([]Object, error) {
	return c.get(unix.RTM_GETCHAIN, nil)
}
