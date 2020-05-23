package tc

import (
	"net"
)

// ipToUint32 converts a legacy ip object to its uint32 representative.
// For IPv6 addresses it returns ErrInvalidArg.
func ipToUint32(ip net.IP) (uint32, error) {
	tmp := ip.To4()
	if tmp == nil {
		return 0, ErrInvalidArg
	}
	return nativeEndian.Uint32(tmp), nil
}

// uint32ToIP converts a legacy ip to a net.IP object.
func uint32ToIP(ip uint32) net.IP {
	netIP := make(net.IP, 4)
	nativeEndian.PutUint32(netIP, ip)
	return netIP
}
