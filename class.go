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

// Add creats a new class
func (c *Class) Add() error {
	return ErrNotImplemented
}

// Replace add/remove a class. If the node does not exist yet it is created
func (c *Class) Replace() error {
	return ErrNotImplemented
}

// Delete removes a class
func (c *Class) Delete() error {
	return ErrNotImplemented
}
