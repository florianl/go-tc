package tc

import (
	"bytes"
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	peditIPProtoTCP = 6
	peditIPProtoUDP = 17
)

const (
	tcaPeditUnspec = iota
	tcaPeditTm
	tcaPeditParms
	tcaPeditPad
	tcaPeditParmsEx
	tcaPeditKeysEx
	tcaPeditKeyEx
)

const (
	tcaPeditKeyExHType = 1
	tcaPeditKeyExCmd   = 2
)

type PeditHeaderType uint16

const (
	PeditHeaderTypeNetwork PeditHeaderType = iota
	PeditHeaderTypeEth
	PeditHeaderTypeIP4
	PeditHeaderTypeIP6
	PeditHeaderTypeTCP
	PeditHeaderTypeUDP
)

type PeditCmd uint16

const (
	PeditCmdSet PeditCmd = iota
	PeditCmdAdd
)

// Pedit contains attributes of the pedit action.
type Pedit struct {
	Sel    PeditSel
	Keys   []PeditKey
	KeysEx []PeditKeyEx
	Tm     *Tcft
}

type PeditSel struct {
	Index   uint32
	Capab   uint32
	Action  int32
	RefCnt  int32
	BindCnt int32
	NKeys   uint8
	Flags   uint8
	Pad     [2]byte
}

type PeditKey struct {
	Mask    uint32
	Val     uint32
	Off     uint32
	At      uint32
	OffMask uint32
	Shift   uint32
}

type PeditKeyEx struct {
	HeaderType PeditHeaderType
	Cmd        PeditCmd
}

func (p *Pedit) SetEthDst(mac net.HardwareAddr) {
	if len(mac) < 6 {
		return
	}
	p.addKey(PeditKey{Val: nativeEndian.Uint32(mac[:4])}, PeditHeaderTypeEth, PeditCmdSet)
	p.addKey(PeditKey{
		Val:  uint32(nativeEndian.Uint16(mac[4:])),
		Mask: 0xffff0000,
		Off:  4,
	}, PeditHeaderTypeEth, PeditCmdSet)
}

func (p *Pedit) SetEthSrc(mac net.HardwareAddr) {
	if len(mac) < 6 {
		return
	}
	p.addKey(PeditKey{
		Val:  uint32(nativeEndian.Uint16(mac[:2])) << 16,
		Mask: 0x0000ffff,
		Off:  4,
	}, PeditHeaderTypeEth, PeditCmdSet)
	p.addKey(PeditKey{
		Val: nativeEndian.Uint32(mac[2:]),
		Off: 8,
	}, PeditHeaderTypeEth, PeditCmdSet)
}

func (p *Pedit) SetIPv4Src(ip net.IP) {
	ip4 := ip.To4()
	if ip4 == nil {
		return
	}
	p.addKey(PeditKey{
		Val: nativeEndian.Uint32(ip4),
		Off: 12,
	}, PeditHeaderTypeIP4, PeditCmdSet)
}

func (p *Pedit) SetIPv4Dst(ip net.IP) {
	ip4 := ip.To4()
	if ip4 == nil {
		return
	}
	p.addKey(PeditKey{
		Val: nativeEndian.Uint32(ip4),
		Off: 16,
	}, PeditHeaderTypeIP4, PeditCmdSet)
}

func (p *Pedit) SetUDPSrc(port uint16) {
	p.SetSrcPort(port, peditIPProtoUDP)
}

func (p *Pedit) SetUDPDst(port uint16) {
	p.SetDstPort(port, peditIPProtoUDP)
}

func (p *Pedit) SetSrcPort(port uint16, protocol uint8) {
	headerType, ok := peditProtocolHeaderType(protocol)
	if !ok {
		return
	}
	p.addKey(PeditKey{
		Val:  uint32(peditSwap16(port)),
		Mask: 0xffff0000,
	}, headerType, PeditCmdSet)
}

func (p *Pedit) SetDstPort(port uint16, protocol uint8) {
	headerType, ok := peditProtocolHeaderType(protocol)
	if !ok {
		return
	}
	p.addKey(PeditKey{
		Val:  uint32(peditSwap16(port)) << 16,
		Mask: 0x0000ffff,
	}, headerType, PeditCmdSet)
}

func (p *Pedit) AddIPv4TTL(delta uint8) {
	p.addKey(PeditKey{
		Val:  uint32(delta),
		Mask: 0xffffff00,
		Off:  8,
	}, PeditHeaderTypeIP4, PeditCmdAdd)
}

func (p *Pedit) addKey(key PeditKey, headerType PeditHeaderType, cmd PeditCmd) {
	p.Keys = append(p.Keys, key)
	p.KeysEx = append(p.KeysEx, PeditKeyEx{HeaderType: headerType, Cmd: cmd})
	p.Sel.NKeys = uint8(len(p.Keys))
}

func peditProtocolHeaderType(protocol uint8) (PeditHeaderType, bool) {
	switch protocol {
	case peditIPProtoTCP:
		return PeditHeaderTypeTCP, true
	case peditIPProtoUDP:
		return PeditHeaderTypeUDP, true
	default:
		return PeditHeaderTypeNetwork, false
	}
}

func peditSwap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func marshalPedit(info *Pedit) ([]byte, error) {
	options := []tcOption{}
	if info == nil {
		return []byte{}, fmt.Errorf("Pedit: %w", ErrNoArg)
	}
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if len(info.Keys) != len(info.KeysEx) {
		return []byte{}, fmt.Errorf("pedit keys and extended keys length mismatch")
	}
	if len(info.Keys) > 255 {
		return []byte{}, fmt.Errorf("pedit supports at most 255 keys")
	}

	sel := info.Sel
	sel.NKeys = uint8(len(info.Keys))
	parms, err := marshalStruct(&sel)
	if err != nil {
		return []byte{}, err
	}
	buf := bytes.NewBuffer(parms)
	for i := range info.Keys {
		keyData, err := marshalStruct(&info.Keys[i])
		if err != nil {
			return []byte{}, err
		}
		buf.Write(keyData)
	}
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaPeditParmsEx, Data: buf.Bytes()})

	keyOptions := make([]tcOption, 0, len(info.KeysEx))
	for i := range info.KeysEx {
		keyExData, err := marshalAttributes([]tcOption{
			{Interpretation: vtUint16, Type: tcaPeditKeyExHType, Data: uint16(info.KeysEx[i].HeaderType)},
			{Interpretation: vtUint16, Type: tcaPeditKeyExCmd, Data: uint16(info.KeysEx[i].Cmd)},
		})
		if err != nil {
			return []byte{}, err
		}
		keyOptions = append(keyOptions, tcOption{Interpretation: vtBytes, Type: tcaPeditKeyEx | nlaFNnested, Data: keyExData})
	}
	keysExData, err := marshalAttributes(keyOptions)
	if err != nil {
		return []byte{}, err
	}
	options = append(options, tcOption{Interpretation: vtBytes, Type: tcaPeditKeysEx | nlaFNnested, Data: keysExData})

	return marshalAttributes(options)
}

func unmarshalPedit(data []byte, info *Pedit) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaPeditParmsEx, tcaPeditParms:
			payload := ad.Bytes()
			if len(payload) < 24 {
				return fmt.Errorf("pedit parms too short: %d", len(payload))
			}
			sel := PeditSel{}
			err = unmarshalStruct(payload[:24], &sel)
			multiError = concatError(multiError, err)
			info.Sel = sel
			offset := 24
			for i := uint8(0); i < sel.NKeys && offset+24 <= len(payload); i++ {
				key := PeditKey{}
				err = unmarshalStruct(payload[offset:offset+24], &key)
				multiError = concatError(multiError, err)
				info.Keys = append(info.Keys, key)
				offset += 24
			}
		case tcaPeditTm:
			tcft := &Tcft{}
			err = unmarshalStruct(ad.Bytes(), tcft)
			multiError = concatError(multiError, err)
			info.Tm = tcft
		case tcaPeditKeysEx, tcaPeditPad:
			// The extended key data is only needed by callers that inspect raw netlink attrs.
		default:
			return fmt.Errorf("unmarshalPedit()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
