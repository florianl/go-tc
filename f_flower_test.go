package tc

import (
	"errors"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlower(t *testing.T) {
	actions := []*Action{
		{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 0x1, Capab: 0x0, Action: 0x4, RefCnt: 0x1, BindCnt: 0x1, Eaction: 0x1, IfIndex: 0x2}}},
	}

	tests := map[string]struct {
		val  Flower
		err1 error
		err2 error
	}{
		"simple": {val: Flower{ClassID: uint32Ptr(42)}},
		"allArguments": {val: Flower{
			ClassID:              uint32Ptr(1),
			Indev:                stringPtr("foo"),
			Actions:              &actions,
			KeyEthDst:            netHardwareAddrPtr(net.HardwareAddr([]byte("00:00:5e:00:53:01"))),
			KeyEthDstMask:        netHardwareAddrPtr(net.HardwareAddr([]byte("00:01:5e:00:53:02"))),
			KeyEthSrc:            netHardwareAddrPtr(net.HardwareAddr([]byte("00:02:5e:00:53:03"))),
			KeyEthSrcMask:        netHardwareAddrPtr(net.HardwareAddr([]byte("00:03:5e:00:53:04"))),
			KeyEthType:           uint16Ptr(2),
			KeyIPProto:           uint8Ptr(3),
			KeyIPv4Src:           netIPPtr(net.ParseIP("1.2.3.4")),
			KeyIPv4SrcMask:       netIPPtr(net.ParseIP("255.255.255.0")),
			KeyIPv4Dst:           netIPPtr(net.ParseIP("4.3.2.1")),
			KeyIPv4DstMask:       netIPPtr(net.ParseIP("255.255.0.0")),
			KeyTCPSrc:            uint16Ptr(4),
			KeyTCPDst:            uint16Ptr(5),
			KeyUDPSrc:            uint16Ptr(6),
			KeyUDPDst:            uint16Ptr(7),
			KeyVlanID:            uint16Ptr(8),
			KeyVlanPrio:          uint8Ptr(9),
			KeyVlanEthType:       uint16Ptr(10),
			KeyEncKeyID:          uint32Ptr(11),
			KeyEncIPv4Src:        netIPPtr(net.ParseIP("3.4.1.2")),
			KeyEncIPv4SrcMask:    netIPPtr(net.ParseIP("255.0.0.0")),
			KeyEncIPv4Dst:        netIPPtr(net.ParseIP("4.3.2.1")),
			KeyEncIPv4DstMask:    netIPPtr(net.ParseIP("0.0.0.0")),
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
			InHwCount:            uint32Ptr(53),
			Flags:                uint32Ptr(54),
			KeyPortSrcMin:        uint16Ptr(55),
			KeyPortSrcMax:        uint16Ptr(56),
			KeyPortDstMin:        uint16Ptr(57),
			KeyPortDstMax:        uint16Ptr(58),
			KeyCtState:           uint16Ptr(59),
			KeyCtStateMask:       uint16Ptr(60),
			KeyCtZone:            uint16Ptr(61),
			KeyCtZoneMask:        uint16Ptr(62),
			KeyCtMark:            uint32Ptr(63),
			KeyCtMarkMask:        uint32Ptr(64),
			KeyHash:              uint32Ptr(65),
			KeyHashMask:          uint32Ptr(66),
			KeyNumOfVLANS:        uint8Ptr(67),
			KeyPppoeSID:          uint16Ptr(68),
			KeyPppProto:          uint16Ptr(69),
			KeyL2TPV3SID:         uint32Ptr(70),
			L2Miss:               uint8Ptr(71),
			KeySpi:               uint32Ptr(72),
			KeySpiMask:           uint32Ptr(73),
			KeyEncFlags:          uint32Ptr(74),
			KeyEncFlagsMask:      uint32Ptr(75),
		}},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlower(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Flower{}
			err2 := unmarshalFlower(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
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
