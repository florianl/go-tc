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

// Add creates a new chain
func (c *Chain) Add(info *Object) error {
	return ErrNotImplemented
}

// Delete removes a chain
func (c *Chain) Delete() error {
	return ErrNotImplemented
}

// Get fetches chains
func (c *Chain) Get() ([]Object, error) {
	return c.get(unix.RTM_GETCHAIN, nil)
}
