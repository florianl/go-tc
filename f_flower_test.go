package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlower(t *testing.T) {
	tests := map[string]struct {
		val  Flower
		err1 error
		err2 error
	}{
		"simple": {val: Flower{ClassID: uint32Ptr(42)}},
		"allAguments": {val: Flower{
			ClassID:              uint32Ptr(1),
			Indev:                stringPtr("foo"),
			KeyEthType:           uint16Ptr(2),
			KeyIPProto:           uint8Ptr(3),
			KeyIPv4Src:           netIPPtr(net.ParseIP("1.1.1.1")),
			KeyIPv4SrcMask:       netIPPtr(net.ParseIP("255.255.255.255")),
			KeyIPv4Dst:           netIPPtr(net.ParseIP("2.2.2.2")),
			KeyIPv4DstMask:       netIPPtr(net.ParseIP("255.255.255.0")),
			KeyTCPSrc:            uint16Ptr(4),
			KeyTCPDst:            uint16Ptr(5),
			KeyUDPSrc:            uint16Ptr(6),
			KeyUDPDst:            uint16Ptr(7),
			KeyVlanID:            uint16Ptr(8),
			KeyVlanPrio:          uint8Ptr(9),
			KeyVlanEthType:       uint16Ptr(10),
			KeyEncKeyID:          uint32Ptr(11),
			KeyEncIPv4Src:        netIPPtr(net.ParseIP("3.3.3.3")),
			KeyEncIPv4SrcMask:    netIPPtr(net.ParseIP("255.255.0.0")),
			KeyEncIPv4Dst:        netIPPtr(net.ParseIP("4.4.4.4")),
			KeyEncIPv4DstMask:    netIPPtr(net.ParseIP("255.0.0.0")),
			KeyTCPSrcMask:        uint16Ptr(12),
			KeyTCPDstMask:        uint16Ptr(13),
			KeyUDPSrcMask:        uint16Ptr(14),
			KeyUDPDstMask:        uint16Ptr(15),
			KeySctpSrc:           uint16Ptr(16),
			KeySctpDst:           uint16Ptr(17),
			KeyEncUDPSrcPort:     uint16Ptr(18),
			KeyEncUDPSrcPortMask: uint16Ptr(19),
			KeyEncUDPDstPort:     uint16Ptr(20),
			KeyEncUDPDstPortMask: uint16Ptr(21),
			KeyFlags:             uint32Ptr(22),
			KeyFlagsMask:         uint32Ptr(23),
			KeyIcmpv4Code:        uint8Ptr(24),
			KeyIcmpv4CodeMask:    uint8Ptr(25),
			KeyIcmpv4Type:        uint8Ptr(26),
			KeyIcmpv4TypeMask:    uint8Ptr(27),
			KeyIcmpv6Code:        uint8Ptr(28),
			KeyIcmpv6CodeMask:    uint8Ptr(29),
			KeyArpSIP:            uint32Ptr(30),
			KeyArpSIPMask:        uint32Ptr(31),
			KeyArpTIP:            uint32Ptr(32),
			KeyArpTIPMask:        uint32Ptr(33),
			KeyArpOp:             uint8Ptr(34),
			KeyArpOpMask:         uint8Ptr(35),
			KeyMplsTTL:           uint8Ptr(36),
			KeyMplsBos:           uint8Ptr(37),
			KeyMplsTc:            uint8Ptr(38),
			KeyMplsLabel:         uint32Ptr(39),
			KeyTCPFlags:          uint16Ptr(40),
			KeyTCPFlagsMask:      uint16Ptr(41),
			KeyIPTOS:             uint8Ptr(42),
			KeyIPTOSMask:         uint8Ptr(43),
			KeyIPTTL:             uint8Ptr(44),
			KeyIPTTLMask:         uint8Ptr(45),
			KeyCVlanID:           uint16Ptr(46),
			KeyCVlanPrio:         uint8Ptr(47),
			KeyCVlanEthType:      uint16Ptr(48),
			KeyEncIPTOS:          uint8Ptr(49),
			KeyEncIPTOSMask:      uint8Ptr(50),
			KeyEncIPTTL:          uint8Ptr(51),
			KeyEncIPTTLMask:      uint8Ptr(52),
		}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlower(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Flower{}
			err2 := unmarshalFlower(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Flower missmatch (want +got):\n%s", diff)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		_, err := marshalFlower(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
