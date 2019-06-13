//+build linux

package tc

import "errors"

// Various errors
var (
	ErrNotImplemented = errors.New("Functionality not yet implemented")
	ErrNoArg          = errors.New("Missing argument")
)

// Config contains options for RTNETLINK
type Config struct {
	// NetNS defines the network namespace
	NetNS int
}

// Constants to define the direction
const (
	// Ingress and Egress can be used as value in Msg.Parent
	Ingress = 0xFFFFFFF2
	Egress  = 0xFFFFFFF3
	// To alter filter in shared blocks, set Msg.Ifindex to MagicBlock
	MagicBlock = 0xFFFFFFFF
)
