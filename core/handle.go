package core

// constants from include/uapi/linux/pkt_sched.h
const (
	handleMajMask uint32 = 0xFFFF0000
	handleMinMask uint32 = 0x0000FFFF
)

// BuildHandle is a simple helper function to construct the handle for the Tcmsg struct
func BuildHandle(maj, min uint32) uint32 {
	return (((maj << 16) & handleMajMask) | (min & handleMinMask))
}

// SplitHandle extracts the major and minor part from a given handle
func SplitHandle(handle uint32) (major, minor uint32) {
	major = (handle & handleMajMask) >> 16
	minor = handle & handleMinMask
	return major, minor
}

// FilterInfo is a simple helper to set the priority and protocol of the filter
func FilterInfo(priority uint16, protocol uint16) uint32 {
	return BuildHandle(uint32(priority), uint32(endianSwapUint16(protocol)))
}

func endianSwapUint16(in uint16) uint16 {
	return (in << 8) | (in >> 8)
}
