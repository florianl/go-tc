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

// ipMaskToUint32 converts a legacy IPMask object to its uint32 representative.
// For IPv6 masks it returns ErrInvalidArg.
func ipMaskToUint32(ipv4 net.IPMask) (uint32, error) {
	ones, bits := ipv4.Size()
	if bits != 32 {
		return 0, ErrInvalidArg
	}
	return uint32(ones), nil
}

// uint32ToIPMask converts a legacy ip mask to a net.IPMask object.
func uint32ToIPMask(bits uint32) net.IPMask {
	return net.CIDRMask(int(bits), 32)
}

// bytesToIP converts a slice of bytes into a net.IP object.
func bytesToIP(ip []byte) (net.IP, error) {
	if len(ip) != net.IPv4len && len(ip) != net.IPv6len {
		return nil, ErrInvalidArg
	}
	return net.IP(ip), nil
}

// ipToBytes casts a ip object into its byte slice representative.
func ipToBytes(ip net.IP) []byte {
	return []byte(ip)
}

// bytesToHardwareAddr converts a slice of bytes into a net.HardwareAddr object.
func bytesToHardwareAddr(mac []byte) net.HardwareAddr {
	return net.HardwareAddr(mac[:])
}

// hardwareAddrToBytes casts a net.HardwareAddr object into its byte slice representative.
func hardwareAddrToBytes(mac net.HardwareAddr) []byte {
	return []byte(mac)
}
