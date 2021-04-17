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
			KeyIPv4SrcMask:       netIPMaskPtr(net.CIDRMask(20, 32)),
			KeyIPv4Dst:           netIPPtr(net.ParseIP("4.3.2.1")),
			KeyIPv4DstMask:       netIPMaskPtr(net.CIDRMask(21, 32)),
			KeyTCPSrc:            uint16Ptr(4),
			KeyTCPDst:            uint16Ptr(5),
			KeyUDPSrc:            uint16Ptr(6),
			KeyUDPDst:            uint16Ptr(7),
			KeyVlanID:            uint16Ptr(8),
			KeyVlanPrio:          uint8Ptr(9),
			KeyVlanEthType:       uint16Ptr(10),
			KeyEncKeyID:          uint32Ptr(11),
			KeyEncIPv4Src:        netIPPtr(net.ParseIP("3.4.1.2")),
			KeyEncIPv4SrcMask:    netIPMaskPtr(net.CIDRMask(22, 32)),
			KeyEncIPv4Dst:        netIPPtr(net.ParseIP("4.3.2.1")),
			KeyEncIPv4DstMask:    netIPMaskPtr(net.CIDRMask(23, 32)),
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
		}},
	}

	endianessMix := make(map[uint16]valueType)
	endianessMix[tcaFlowerKeyEthType] = vtUint16Be
	endianessMix[tcaFlowerKeyIPv4Src] = vtUint32Be
	endianessMix[tcaFlowerKeyIPv4SrcMask] = vtUint32Be
	endianessMix[tcaFlowerKeyIPv4Dst] = vtUint32Be
	endianessMix[tcaFlowerKeyIPv4DstMask] = vtUint32Be
	endianessMix[tcaFlowerKeyTCPSrc] = vtUint16Be
	endianessMix[tcaFlowerKeyTCPDst] = vtUint16Be
	endianessMix[tcaFlowerKeyUDPSrc] = vtUint16Be
	endianessMix[tcaFlowerKeyUDPDst] = vtUint16Be
	endianessMix[tcaFlowerKeyVlanEthType] = vtUint16Be
	endianessMix[tcaFlowerKeyEncKeyID] = vtUint32Be
	endianessMix[tcaFlowerKeyEncIPv4Src] = vtUint32Be
	endianessMix[tcaFlowerKeyEncIPv4SrcMask] = vtUint32Be
	endianessMix[tcaFlowerKeyEncIPv4Dst] = vtUint32Be
	endianessMix[tcaFlowerKeyEncIPv4DstMask] = vtUint32Be
	endianessMix[tcaFlowerKeyTCPSrcMask] = vtUint16Be
	endianessMix[tcaFlowerKeyTCPDstMask] = vtUint16Be
	endianessMix[tcaFlowerKeyUDPSrcMask] = vtUint16Be
	endianessMix[tcaFlowerKeyUDPDstMask] = vtUint16Be
	endianessMix[tcaFlowerKeySCTPSrcMask] = vtUint16Be
	endianessMix[tcaFlowerKeySCTPDstMask] = vtUint16Be
	endianessMix[tcaFlowerKeySCTPSrc] = vtUint16Be
	endianessMix[tcaFlowerKeySCTPDst] = vtUint16Be
	endianessMix[tcaFlowerKeyEncUDPSrcPort] = vtUint16Be
	endianessMix[tcaFlowerKeyEncUDPSrcPortMask] = vtUint16Be
	endianessMix[tcaFlowerKeyEncUDPDstPort] = vtUint16Be
	endianessMix[tcaFlowerKeyEncUDPDstPortMask] = vtUint16Be
	endianessMix[tcaFlowerKeyFlags] = vtUint32Be
	endianessMix[tcaFlowerKeyFlagsMask] = vtUint32Be
	endianessMix[tcaFlowerKeyArpSIP] = vtUint32Be
	endianessMix[tcaFlowerKeyArpSIPMask] = vtUint32Be
	endianessMix[tcaFlowerKeyArpTIP] = vtUint32Be
	endianessMix[tcaFlowerKeyArpTIPMask] = vtUint32Be
	endianessMix[tcaFlowerKeyTCPFlags] = vtUint16Be
	endianessMix[tcaFlowerKeyTCPFlagsMask] = vtUint16Be
	endianessMix[tcaFlowerKeyCVlanEthType] = vtUint16Be

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalFlower(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			newData := changeEndianess(t, data, endianessMix)

			val := Flower{}
			err2 := unmarshalFlower(newData, &val)
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
