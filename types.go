package tc

import "errors"

// Various errors
var (
	ErrNotImplemented = errors.New("functionality not yet implemented")
	ErrNoArg          = errors.New("missing argument")
	ErrNoArgAlter     = errors.New("argument cannot be altered")
	ErrNotLinux       = errors.New("not implemented for OS other than linux")
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
