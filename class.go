//+build linux

package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

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
func (c *Class) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateClassObject(rtmNewClass, info)
	if err != nil {
		return err
	}
	return c.action(rtmNewClass, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a class. If the node does not exist yet it is created
func (c *Class) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateClassObject(rtmNewClass, info)
	if err != nil {
		return err
	}
	return c.action(rtmNewClass, netlink.Create, &info.Msg, options)
}

// Delete removes a class
func (c *Class) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateClassObject(rtmDelClass, info)
	if err != nil {
		return err
	}
	return c.action(rtmDelClass, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all classes
func (c *Class) Get(i *Msg) ([]Object, error) {
	if i == nil {
		return []Object{}, ErrNoArg
	}
	return c.get(rtmGetClass, i)
}

func validateClassObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, fmt.Errorf("Could not set device ID 0")
	}
	return options, nil
}
