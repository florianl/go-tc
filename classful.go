//+build linux

package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaQfqUnspec = iota
	tcaQfqWeight
	tcaQfqLmax
)

// Qfq contains attributes of the qfq discipline
type Qfq struct {
	Weight uint32
	Lmax   uint32
}

func extractQfqOptions(data []byte, info *Qfq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaQfqWeight:
			info.Weight = ad.Uint32()
		case tcaQfqLmax:
			info.Lmax = ad.Uint32()
		default:
			return fmt.Errorf("extractQfqOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaHfscUnspec = iota
	tcaHfscRsc
	tcaHfscFsc
	tcaHfscUsc
)

// Hfsc contains attributes of the hfsc discipline
type Hfsc struct {
	Rsc *ServiceCurve
	Fsc *ServiceCurve
	Usc *ServiceCurve
}

func extractHfscOptions(data []byte, info *Hfsc) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHfscRsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Rsc = curve
		case tcaHfscFsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Fsc = curve
		case tcaHfscUsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Usc = curve
		default:
			return fmt.Errorf("extractHfscOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// ServiceCurve from include/uapi/linux/pkt_sched.h
type ServiceCurve struct {
	M1 uint32
	D  uint32
	M2 uint32
}

func extractServiceCurve(data []byte, info *ServiceCurve) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaHtbUnspec = iota
	tcaHtbParms
	tcaHtbInit
	tcaHtbCtab
	tcaHtbRtab
	tcaHtbDirectQlen
	tcaHtbRate64
	tcaHtbCeil64
	tcaHtbPad
)

// Htb contains attributes of the HTB discipline
type Htb struct {
	Parms      *HtbOpt
	Init       *HtbGlob
	Ctab       []byte
	Rtab       []byte
	DirectQlen uint32
	Rate64     uint64
	Ceil64     uint64
}

func extractHtbOptions(data []byte, info *Htb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHtbParms:
			opt := &HtbOpt{}
			if err := extractHtbOpt(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		case tcaHtbInit:
			glob := &HtbGlob{}
			if err := extractHtbGlob(ad.Bytes(), glob); err != nil {
				return err
			}
			info.Init = glob
		case tcaHtbCtab:
			info.Ctab = ad.Bytes()
		case tcaHtbRtab:
			info.Rtab = ad.Bytes()
		case tcaHtbDirectQlen:
			info.DirectQlen = ad.Uint32()
		case tcaHtbRate64:
			info.Rate64 = ad.Uint64()
		case tcaHtbCeil64:
			info.Ceil64 = ad.Uint64()
		case tcaHtbPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("extractHtbOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// HtbGlob from include/uapi/linux/pkt_sched.h
type HtbGlob struct {
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
}

func extractHtbGlob(data []byte, info *HtbGlob) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// HtbOpt from include/uapi/linux/pkt_sched.h
type HtbOpt struct {
	Rate    RateSpec
	Ceil    RateSpec
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}

func extractHtbOpt(data []byte, info *HtbOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaDsmarkUnspec = iota
	tcaDsmarkIndices
	tcaDsmarkDefaultIndex
	tcaDsmarkSetTCIndex
	tcaDsmarkMask
	tcaDsmarkValue
)

// Dsmark contains attributes of the dsmark discipline
type Dsmark struct {
	Indices      uint16
	DefaultIndex uint16
	// SetTCIndex NLA_FLAG
	Mask  uint8
	Value uint8
}

func extractDsmarkOptions(data []byte, info *Dsmark) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaDsmarkIndices:
			info.Indices = ad.Uint16()
		case tcaDsmarkDefaultIndex:
			info.DefaultIndex = ad.Uint16()
		case tcaDsmarkSetTCIndex:
			// TODO: NLA_FLAG not yet supported
		case tcaDsmarkMask:
			info.Mask = ad.Uint8()
		case tcaDsmarkValue:
			info.Value = ad.Uint8()
		default:
			return fmt.Errorf("extractDsmarkOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaDrrUnspec = iota
	tcaDrrQuantum
)

// Drr contains attributes of the drr discipline
type Drr struct {
	Quantum uint32
}

func extractDrrOptions(data []byte, info *Drr) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaDrrQuantum:
			info.Quantum = ad.Uint32()
		default:
			return fmt.Errorf("extractDrrOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaCbqUnspec = iota
	tcaCbqLssOpt
	tcaCbqWrrOpt
	tcaCbqFOpt
	tcaCbqOVLStrategy
	tcaCbqRate
	tcaCbqRTab
	tcaCbqPolice
)

// Cbq contains attributes of the cbq discipline
type Cbq struct {
	LssOpt      *CbqLssOpt
	WrrOpt      *CbqWrrOpt
	FOpt        *CbqFOpt
	OVLStrategy *CbqOvl
	Rate        *RateSpec
	RTab        []byte
	Police      *CbqPolice
}

func extractCbqOptions(data []byte, info *Cbq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaCbqLssOpt:
			arg := &CbqLssOpt{}
			if err := extractCbqLssOpt(ad.Bytes(), arg); err != nil {
				return err
			}
			info.LssOpt = arg
		case tcaCbqWrrOpt:
			arg := &CbqWrrOpt{}
			if err := extractCbqWrrOpt(ad.Bytes(), arg); err != nil {
				return err
			}
			info.WrrOpt = arg
		case tcaCbqFOpt:
			arg := &CbqFOpt{}
			if err := extractCbqFOpt(ad.Bytes(), arg); err != nil {
				return err
			}
			info.FOpt = arg
		case tcaCbqOVLStrategy:
			arg := &CbqOvl{}
			if err := extractCbqOvl(ad.Bytes(), arg); err != nil {
				return err
			}
			info.OVLStrategy = arg
		case tcaCbqRate:
			arg := &RateSpec{}
			if err := extractRateSpec(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Rate = arg
		case tcaCbqRTab:
			info.RTab = ad.Bytes()
		case tcaCbqPolice:
			arg := &CbqPolice{}
			if err := extractCbqPolice(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Police = arg
		default:
			return fmt.Errorf("extractCbqOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// CbqLssOpt from include/uapi/linux/pkt_sched.h
type CbqLssOpt struct {
	Change  byte
	Flags   byte
	EwmaLog byte
	Level   byte
	Maxidle uint32
	Minidle uint32
	OffTime uint32
	Avpkt   uint32
}

func extractCbqLssOpt(data []byte, info *CbqLssOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// CbqWrrOpt from include/uapi/linux/pkt_sched.h
type CbqWrrOpt struct {
	Flags     byte
	Priority  byte
	CPriority byte
	Reserved  byte
	Allot     uint32
	Weight    uint32
}

func extractCbqWrrOpt(data []byte, info *CbqWrrOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// CbqFOpt from include/uapi/linux/pkt_sched.h
type CbqFOpt struct {
	split     uint32
	defmap    uint32
	defchange uint32
}

func extractCbqFOpt(data []byte, info *CbqFOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// CbqOvl from include/uapi/linux/pkt_sched.h
type CbqOvl struct {
	strategy  byte
	priority2 byte
	pad       uint16
	penalty   uint32
}

func extractCbqOvl(data []byte, info *CbqOvl) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// CbqPolice from include/uapi/linux/pkt_sched.h
type CbqPolice struct {
	police byte
	Res1   byte
	Res2   uint16
}

func extractCbqPolice(data []byte, info *CbqPolice) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaAtmUnspec = iota
	tcaAtmFD
	tcaAtmPtr
	tcaAtmHdr
	tcaAtmExcess
	tcaAtmAddr
	tcaAtmState
)

// Atm contains attributes of the atm discipline
type Atm struct {
	FD     uint32
	Excess uint32
	Addr   *AtmPvc
	State  uint32
}

func extractAtmOptions(data []byte, info *Atm) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaAtmFD:
			info.FD = ad.Uint32()
		case tcaAtmExcess:
			info.Excess = ad.Uint32()
		case tcaAtmAddr:
			arg := &AtmPvc{}
			if err := extractAtmPvc(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Addr = arg
		case tcaAtmState:
			info.State = ad.Uint32()
		default:
			return fmt.Errorf("extractAtmOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return nil
}

// AtmPvc from include/uapi/linux/atm.h
type AtmPvc struct {
	SapFamily byte
	Itf       byte
	Vpi       byte
	Vci       byte
}

func extractAtmPvc(data []byte, info *AtmPvc) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}
