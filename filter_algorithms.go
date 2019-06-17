//+build linux

package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaBpfUnspec = iota
	tcaBpfAct
	tcaBpfPolice
	tcaBpfClassid
	tcaBpfOpsLen
	tcaBpfOps
	tcaBpfFd
	tcaBpfName
	tcaBpfFlags
	tcaBpfFlagsGen
	tcaBpfTag
	tcaBpfID
)

// BPF contains attributes of the bpf discipline
type BPF struct {
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      []byte
	ID       uint32
	Action   *Action
}

func validateBPFOptions(info *BPF) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("BPF options are missing")
	}

	// TODO: improve logic and check combinations

	if info.FD != 0 && len(info.Ops) != 0 {
		return []byte{}, fmt.Errorf("can not use FD and Ops at the same time")
	}

	if info.FD != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFd, Data: info.FD})
	}

	if len(info.Ops) != 0 {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaBpfOpsLen, Data: info.OpsLen})
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfOps, Data: info.Ops})
	}

	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfClassid, Data: info.ClassID})
	}
	if len(info.Name) != 0 {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaBpfName, Data: info.Name})
	}
	if info.Flags != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlags, Data: info.Flags})
	}
	if info.FlagsGen != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlagsGen, Data: info.FlagsGen})
	}

	return marshalAttributes(options)

}

func extractBpfOptions(data []byte, info *BPF) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaBpfAct:
			action := &Action{}
			if err := extractBPFAction(ad.Bytes(), action); err != nil {
				return err
			}
			info.Action = action
		case tcaBpfClassid:
			info.ClassID = ad.Uint32()
		case tcaBpfOpsLen:
			info.OpsLen = ad.Uint16()
		case tcaBpfOps:
			info.Ops = ad.Bytes()
		case tcaBpfFd:
			info.FD = ad.Uint32()
		case tcaBpfName:
			info.Name = ad.String()
		case tcaBpfFlags:
			info.Flags = ad.Uint32()
		case tcaBpfFlagsGen:
			info.FlagsGen = ad.Uint32()
		case tcaBpfTag:
			info.Tag = ad.Bytes()
		case tcaBpfID:
			info.ID = ad.Uint32()
		default:
			return fmt.Errorf("extractBpfOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// Actual attributes are nested inside and the nested bit is not set :-/
func extractBPFAction(data []byte, action *Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	ad.Next()
	return extractTcAction(ad.Bytes(), action)
}

func extractActBpfOptions(data []byte, attr *BPFActionOptions) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaActBpfTm:
			info := &Tcft{}
			if err := extractTcft(ad.Bytes(), info); err != nil {
				return err
			}
			attr.Tcft = info
		case tcaActBpfParms:
			info := &ActBpf{}
			if err := extractTcActBpf(ad.Bytes(), info); err != nil {
				return err
			}
			attr.Act = info
		case tcaActBpfOpsLen:
			attr.OpsLen = ad.Uint16()
		case tcaActBpfOps:
			attr.Ops = ad.Bytes()
		case tcaActBpfFD:
			attr.FD = ad.Uint32()
		case tcaActBpfName:
			attr.Name = ad.String()
		default:
			return fmt.Errorf("extractActBpfOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// BPFActionOptions contains various action attributes
type BPFActionOptions struct {
	OpsLen uint16
	Ops    []byte
	Tcft   *Tcft
	FD     uint32
	Name   string
	Act    *ActBpf
}

const (
	tcaActBpfUnspec = iota
	tcaActBpfTm
	tcaActBpfParms
	tcaActBpfOpsLen
	tcaActBpfOps
	tcaActBpfFD
	tcaActBpfName
	tcaActBpfPad
	tcaActBpfTag
	tcaActBpfID
)

// ActionStats contains various statistics of a action
type ActionStats struct {
	Basic     *GenStatsBasic
	RateEst   *GenStatsRateEst
	Queue     *GenStatsQueue
	RateEst64 *GenStatsRateEst64
}

// Action describes a Traffic Control action
type Action struct {
	Kind       string
	Statistics *ActionStats
	BPFOptions *BPFActionOptions
}

const (
	tcaRsvpUnspec = iota
	tcaRsvpClassID
	tcaRsvpDst
	tcaRsvpSrc
	tcaRsvpPInfo
	tcaRsvpPolice
	tcaRsvpAct
)

// Rsvp contains attributes of the rsvp discipline
type Rsvp struct {
	ClassID uint32
	Dst     []byte
	Src     []byte
	PInfo   *RsvpPInfo
}

func extractRsvpOptions(data []byte, info *Rsvp) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaRsvpClassID:
			info.ClassID = ad.Uint32()
		case tcaRsvpDst:
			info.Dst = ad.Bytes()
		case tcaRsvpSrc:
			info.Src = ad.Bytes()
		case tcaRsvpPInfo:
			arg := &RsvpPInfo{}
			if err := extractRsvpPInfo(ad.Bytes(), arg); err != nil {
				return err
			}
			info.PInfo = arg
		default:
			return fmt.Errorf("extractRsvpOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil

}

// RsvpPInfo from include/uapi/linux/pkt_sched.h
type RsvpPInfo struct {
	Dpi       *RsvpGpi
	Spi       *RsvpGpi
	Protocol  uint8
	TunnelID  uint8
	TunnelHdr uint8
	Pad       uint8
}

func extractRsvpPInfo(data []byte, info *RsvpPInfo) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// RsvpGpi from include/uapi/linux/pkt_sched.h
type RsvpGpi struct {
	Key    uint32
	Mask   uint32
	Offset uint32
}

const (
	tcaRoute4Unspec = iota
	tcaRoute4ClassID
	tcaRoute4To
	tcaRoute4From
	tcaRoute4IIf
	tcaRoute4Police
	tcaRoute4Act
)

// Route4 contains attributes of the route discipline
type Route4 struct {
	ClassID uint32
	To      uint32
	From    uint32
	IIf     uint32
}

func extractRoute4Options(data []byte, info *Route4) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaRoute4ClassID:
			info.ClassID = ad.Uint32()
		case tcaRoute4To:
			info.To = ad.Uint32()
		case tcaRoute4From:
			info.From = ad.Uint32()
		case tcaRoute4IIf:
			info.IIf = ad.Uint32()
		default:
			return fmt.Errorf("extractRoute4Options()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil

}

const (
	tcaFwUnspec = iota
	tcaFwClassID
	tcaFwPolice
	tcaFwInDev
	tcaFwAct
	tcaFwMask
)

// Fw contains attributes of the fw discipline
type Fw struct {
	ClassID uint32
	InDev   string
	Mask    uint32
}

func extractFwOptions(data []byte, info *Fw) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFwClassID:
			info.ClassID = ad.Uint32()
		case tcaFwInDev:
			info.InDev = ad.String()
		case tcaFwMask:
			info.Mask = ad.Uint32()
		default:
			return fmt.Errorf("extractFwOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaFlowUnspec = iota
	tcaFlowKeys
	tcaFlowMode
	tcaFlowBaseClass
	tcaFlowRShift
	tcaFlowAddend
	tcaFlowMask
	tcaFlowXOR
	tcaFlowDivisor
	tcaFlowAct
	tcaFlowPolice
	tcaFlowEMatches
	tcaFlowPerTurb
)

// Flow contains attributes of the flow discipline
type Flow struct {
	Keys      uint32
	Mode      uint32
	BaseClass uint32
	RShift    uint32
	Addend    uint32
	Mask      uint32
	XOR       uint32
	Divisor   uint32
	PerTurb   uint32
}

func extractFlowOptions(data []byte, info *Flow) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFlowKeys:
			info.Keys = ad.Uint32()
		case tcaFlowMode:
			info.Mode = ad.Uint32()
		case tcaFlowBaseClass:
			info.BaseClass = ad.Uint32()
		case tcaFlowRShift:
			info.RShift = ad.Uint32()
		case tcaFlowAddend:
			info.Addend = ad.Uint32()
		case tcaFlowMask:
			info.Mask = ad.Uint32()
		case tcaFlowXOR:
			info.XOR = ad.Uint32()
		case tcaFlowDivisor:
			info.Divisor = ad.Uint32()
		case tcaFlowPerTurb:
			info.PerTurb = ad.Uint32()
		default:
			return fmt.Errorf("extractFlowOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
