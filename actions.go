package tc

import (
	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Actions represents the actions part of rtnetlink
type Actions struct {
	Tc
}

// tcamsg is Actions specific
type tcaMsg struct {
	family uint8
	_      uint8  // pad1
	_      uint16 // pad2
}

// Actions allows to read and alter actions
func (tc *Tc) Actions() *Actions {
	return &Actions{*tc}
}

// Add creates a new actions
func (a *Actions) Add(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_NEWACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_NEWACTION, netlink.Create|netlink.Excl, tcaMsg{
		family: unix.AF_UNSPEC,
	}, options)
}

// Replace add/remove an actions. If the node does not exist yet it is created
func (a *Actions) Replace(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_NEWACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_NEWACTION, netlink.Create, tcaMsg{
		family: unix.AF_UNSPEC,
	}, options)
}

// Delete removes an actions
func (a *Actions) Delete(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_DELACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_DELACTION, netlink.HeaderFlags(0), tcaMsg{
		family: unix.AF_UNSPEC,
	}, options)
}

func validateActionsObject(cmd int, info []*Action) ([]tcOption, error) {
	options := []tcOption{}

	data, err := marshalActions(cmd, info)
	if err != nil {
		return options, err
	}
	options = append(options, tcOption{Interpretation: vtBytes, Type: 1 /*TCA_ROOT_TAB*/, Data: data})

	return options, nil
}
