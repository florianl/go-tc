// +build linux

package unix

import linux "golang.org/x/sys/unix"

type IfInfomsg = linux.IfInfomsg

const (
	AF_UNSPEC     = linux.AF_UNSPEC
	NETLINK_ROUTE = linux.NETLINK_ROUTE
	IFLA_EXT_MASK = linux.IFLA_EXT_MASK
	RTM_GETLINK   = linux.RTM_GETLINK
	RTNLGRP_TC    = linux.RTNLGRP_TC
)
