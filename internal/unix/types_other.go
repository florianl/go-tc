// +build !linux

package unix

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

const (
	AF_UNSPEC     = 0x0
	NETLINK_ROUTE = 0x0
	IFLA_EXT_MASK = 0x1d
	RTM_GETLINK   = 0x12
	RTNLGRP_TC    = 0x4
)
