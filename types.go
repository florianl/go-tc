package tc

import "errors"

// Various errors
var (
	ErrNotImplemented = errors.New("Functionallity not yet implemented")
	ErrNoArg          = errors.New("Missing argument")
)

// Config contains options for RTNETLINK
type Config struct {
	// NetNS defines the network namespace
	NetNS int
}

// Constants to define the direction
const (
	Ingress = 0xFFFFFFF2
	Egress  = 0xFFFFFFF3
)
