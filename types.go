package tc

import "errors"

// Various errors
var (
	ErrNotImplemented = errors.New("functionality not yet implemented")
	ErrNoArg          = errors.New("missing argument")
	ErrNoArgAlter     = errors.New("argument cannot be altered")
	ErrNotLinux       = errors.New("not implemented for OS other than linux")
	ErrInvalidDev     = errors.New("invalid device ID")
)

// Config contains options for RTNETLINK
type Config struct {
	// NetNS defines the network namespace
	NetNS int
}

// Constants to define the direction
const (
	HandleRoot    uint32 = 0xFFFFFFFF
	HandleIngress uint32 = 0xFFFFFFF1

	HandleMinPriority uint32 = 0xFFE0
	HandleMinIngress  uint32 = 0xFFF2
	HandleMinEgress   uint32 = 0xFFF3

	// To alter filter in shared blocks, set Msg.Ifindex to MagicBlock
	MagicBlock = 0xFFFFFFFF
)

// constants from include/uapi/linux/pkt_sched.h
const (
	handleMajMask uint32 = 0xFFFF0000
	handleMinMask uint32 = 0x0000FFFF
)
