package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLegacyIPHelper(t *testing.T) {
	tests := map[string]struct {
		input net.IP
		err   error
	}{
		"googleDNSIPv4": {input: net.ParseIP("8.8.8.8")},
		"googleDNSIPv6": {input: net.ParseIP("2001:4860:4860::8888"), err: ErrInvalidArg},
		"invalidIP":     {input: net.ParseIP("foobar"), err: ErrInvalidArg},
	}

	for name, testcase := range tests {
		name := name
		testcase := testcase
		t.Run(name, func(t *testing.T) {
			var num uint32
			var err error
			if num, err = ipToUint32(testcase.input); err != nil {
				if errors.Is(err, testcase.err) {
					t.Log("Received expected error")
					return
				}
				t.Fatalf("Received unexpected error: %v", err)
			}
			if testcase.err != nil {
				t.Fatalf("Expected error but got none")
			}
			got := uint32ToIP(num)
			if diff := cmp.Diff(testcase.input, got); diff != "" {
				t.Fatalf("TestLegacyIPHelper() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRealIPHelper(t *testing.T) {
	tests := map[string]struct {
		input net.IP
		err   error
	}{
		"googleDNSIPv4": {input: net.ParseIP("8.8.8.8")},
		"googleDNSIPv6": {input: net.ParseIP("2001:4860:4860::8888")},
		"invalidIP":     {input: net.ParseIP("foobar"), err: ErrInvalidArg},
	}

	for name, testcase := range tests {
		name := name
		testcase := testcase
		t.Run(name, func(t *testing.T) {
			var ip net.IP
			var err error
			slice := ipToBytes(testcase.input)

			if ip, err = bytesToIP(slice); err != nil {
				if errors.Is(err, testcase.err) {
					t.Log("Received expected error")
					return
				}
				t.Fatalf("Received unexpected error: %v", err)
			}
			if testcase.err != nil {
				t.Fatalf("Expected error but got none")
			}
			if diff := cmp.Diff(testcase.input, ip); diff != "" {
				t.Fatalf("TestRealIPHelper() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("invalid byte length", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			var slice []byte
			for j := 0; j < i; j++ {
				slice = append(slice, byte(j))
			}
			_, err := bytesToIP(slice)
			if err != nil {
				if errors.Is(err, ErrInvalidArg) && (i != net.IPv4len && i != net.IPv6len) {
					t.Logf("Received expected error for byte length of %d", i)
					continue
				}
				t.Fatalf("Received unexpected error for slice length of %d: %v", i, err)
			}
		}
	})
}

func TestMacHelper(t *testing.T) {
	for _, macStr := range []string{
		"00:00:5e:00:53:01",
		"02:00:5e:10:00:00:00:01",
		"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01",
		"00-00-5e-00-53-01",
		"02-00-5e-10-00-00-00-01",
		"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01",
		"0000.5e00.5301",
		"0200.5e10.0000.0001",
		"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001",
	} {
		macStr := macStr
		t.Run(macStr, func(t *testing.T) {
			mac, err := net.ParseMAC(macStr)
			if err != nil {
				t.Fatalf("failed to parse mac string: %v", err)
			}
			tmp := hardwareAddrToBytes(mac)
			mac2 := bytesToHardwareAddr(tmp)
			if diff := cmp.Diff(mac, mac2); diff != "" {
				t.Fatalf("HardwareAddr missmatch (-want +got):\n%s", diff)
			}
		})
	}
}
