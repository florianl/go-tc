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
	HandleRoot    = 0xFFFFFFFF
	HandleIngress = 0xFFFFFFF1

	HandleMinPriority = 0xFFE0
	HandleMinIngress  = 0xFFF2
	HandleMinEgress   = 0xFFF3

	// To alter filter in shared blocks, set Msg.Ifindex to MagicBlock
	MagicBlock = 0xFFFFFFFF
)
