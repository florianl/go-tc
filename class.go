package tc

import (
	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Class represents the class part of rtnetlink
type Class struct {
	Tc
}

// Class allows to read and alter classes
func (tc *Tc) Class() *Class {
	return &Class{*tc}
}

// Add creats a new class
func (c *Class) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWTCLASS, info)
	if err != nil {
		return err
	}
	return c.action(unix.RTM_NEWTCLASS, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a class. If the node does not exist yet it is created
func (c *Class) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWTCLASS, info)
	if err != nil {
		return err
	}
	return c.action(unix.RTM_NEWTCLASS, netlink.Create, &info.Msg, options)
}

// Delete removes a class
func (c *Class) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_DELTCLASS, info)
	if err != nil {
		return err
	}
	return c.action(unix.RTM_DELTCLASS, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all classes
func (c *Class) Get(i *Msg) ([]Object, error) {
	if i == nil {
		return []Object{}, ErrNoArg
	}
	return c.get(unix.RTM_GETTCLASS, i)
}
