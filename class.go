package tc

// Class represents the class part of rtnetlink
type Class struct {
	Tc
}

const (
	rtmNewClass = 40
	rtmDelClass = 41
	rtmGetClass = 42
)

// Class allows to read and alter classes
func (tc *Tc) Class() *Class {
	return &Class{*tc}
}

// New adds a class
func (c *Class) New() error {
	return ErrNotImplemented
}

// Del removes a class
func (c *Class) Del() error {
	return ErrNotImplemented
}

// Get fetches all classes
func (c *Class) Get() error {
	return ErrNotImplemented
}
