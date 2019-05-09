package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaFqCodelUnspec = iota
	tcaFqCodelTarget
	tcaFqCodelLimit
	tcaFqCodelInterval
	tcaFqCodelEcn
	tcaFqCodelFlows
	tcaFqCodelQuantum
	tcaFqCodelCeThreshold
	tcaFqCodelDropBatchSize
	tcaFqCodelMemoryLimit
)

// FqCodel contains attributes of the fq_codel discipline
type FqCodel struct {
	Target        uint32
	Limit         uint32
	Interval      uint32
	ECN           uint32
	Flows         uint32
	Quantum       uint32
	CEThreshold   uint32
	DropBatchSize uint32
	MemoryLimit   uint32
}

func validateFqCodelOptions(info *FqCodel) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("fq_codel options are missing")
	}

	if info.Target != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelTarget, Data: info.Target})
	}

	if info.Limit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelLimit, Data: info.Limit})
	}

	if info.Interval != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelInterval, Data: info.Interval})
	}

	if info.ECN != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelEcn, Data: info.ECN})
	}

	if info.Flows != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelFlows, Data: info.Flows})
	}

	if info.Quantum != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelQuantum, Data: info.Quantum})
	}

	if info.CEThreshold != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelCeThreshold, Data: info.CEThreshold})
	}

	if info.DropBatchSize != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelDropBatchSize, Data: info.DropBatchSize})
	}

	if info.MemoryLimit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCodelMemoryLimit, Data: info.MemoryLimit})
	}

	return marshalAttributes(options)
}

func extractFqCodelOptions(data []byte, info *FqCodel) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFqCodelTarget:
			info.Target = ad.Uint32()
		case tcaFqCodelLimit:
			info.Limit = ad.Uint32()
		case tcaFqCodelInterval:
			info.Interval = ad.Uint32()
		case tcaFqCodelEcn:
			info.ECN = ad.Uint32()
		case tcaFqCodelFlows:
			info.Flows = ad.Uint32()
		case tcaFqCodelQuantum:
			info.Quantum = ad.Uint32()
		case tcaFqCodelCeThreshold:
			info.CEThreshold = ad.Uint32()
		case tcaFqCodelDropBatchSize:
			info.DropBatchSize = ad.Uint32()
		case tcaFqCodelMemoryLimit:
			info.MemoryLimit = ad.Uint32()
		default:
			return fmt.Errorf("extractFqCodelOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaTbfUnspec = iota
	tcaTbfParms
	tcaTbfRtab
	tcaTbfPtab
	tcaTbfRate64
	tcaTbfPrate64
	tcaTbfBurst
	tcaTbfPburst
	tcaTbfPad
)

// Tbf contains attributes of the TBF discipline
type Tbf struct {
	Parms   *TbfQopt
	Rtab    []byte
	Ptab    []byte
	Rate64  uint64
	Prate64 uint64
	Burst   uint32
	Pburst  uint32
}

func extractTbfOptions(data []byte, info *Tbf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaTbfParms:
			qopt := &TbfQopt{}
			if err := extractTbfQopt(ad.Bytes(), qopt); err != nil {
				return err
			}
			info.Parms = qopt
		case tcaTbfRtab:
			info.Rtab = ad.Bytes()
		case tcaTbfPtab:
			info.Ptab = ad.Bytes()
		case tcaTbfRate64:
			info.Rate64 = ad.Uint64()
		case tcaTbfPrate64:
			info.Prate64 = ad.Uint64()
		case tcaTbfBurst:
			info.Burst = ad.Uint32()
		case tcaTbfPburst:
			info.Pburst = ad.Uint32()
		case tcaTbfPad:
			// padding does not contail data, we just skip it
		default:
			return fmt.Errorf("extractTbfOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// TbfQopt from include/uapi/linux/pkt_sched.h
type TbfQopt struct {
	Rate     RateSpec
	PeakRate RateSpec
	Limit    uint32
	Buffer   uint32
	Mtu      uint32
}

func extractTbfQopt(data []byte, info *TbfQopt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaSfbUnspec = iota
	tcaSfbParms
)

// Sfb contains attributes of the SBF discipline
type Sfb struct {
	Parms *SfbQopt
}

func extractSfbOptions(data []byte, info *Sfb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaSfbParms:
			opt := &SfbQopt{}
			if err := extractSfbQopt(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		default:
			return fmt.Errorf("extractSbfOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// SfbQopt from include/uapi/linux/pkt_sched.h
type SfbQopt struct {
	RehashInterval uint32 // in ms
	WarmupTime     uint32 //  in ms
	Max            uint32
	BinSize        uint32
	Increment      uint32
	Decrement      uint32
	Limit          uint32
	PenaltyRate    uint32
	PenaltyBurst   uint32
}

func extractSfbQopt(data []byte, info *SfbQopt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaRedUnspec = iota
	tcaRedParms
	tcaRedStab
	tcaRedMaxP
)

// Red contains attributes of the red discipline
type Red struct {
	Parms *RedQOpt
	MaxP  uint32
}

func extractRedOptions(data []byte, info *Red) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaRedParms:
			opt := &RedQOpt{}
			if err := extractRedQOpt(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		case tcaRedMaxP:
			info.MaxP = ad.Uint32()
		default:
			return fmt.Errorf("extractRedOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// RedQOpt from include/uapi/linux/pkt_sched.h
type RedQOpt struct {
	Limit    uint32
	QthMin   uint32
	QthMax   uint32
	Wlog     byte
	Plog     byte
	ScellLog byte
	Flags    byte
}

func extractRedQOpt(data []byte, info *RedQOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

const (
	tcaMqPrioUnspec = iota
	tcaMqPrioMode
	tcaMqPrioShaper
	tcaMqPrioMinRate64
	tcaMqPrioMaxRate64
)

// MqPrio contains attributes of the mqprio discipline
type MqPrio struct {
	Mode      uint16
	Shaper    uint16
	MinRate64 uint64
	MaxRate64 uint64
}

func extractMqPrioOptions(data []byte, info *MqPrio) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaMqPrioMaxRate64:
			info.Mode = ad.Uint16()
		case tcaMqPrioShaper:
			info.Mode = ad.Uint16()
		default:
			return fmt.Errorf("extractMqPrioOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaCodelUnspec = iota
	tcaCodelTarget
	tcaCodelLimit
	tcaCodelInterval
	tcaCodelECN
	tcaCodelCEThreshold
)

// Codel contains attributes of the codel discipline
type Codel struct {
	Target      uint32
	Limit       uint32
	Interval    uint32
	ECN         uint32
	CEThreshold uint32
}

func extractCodelOptions(data []byte, info *Codel) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaCodelTarget:
			info.Target = ad.Uint32()
		case tcaCodelLimit:
			info.Limit = ad.Uint32()
		case tcaCodelInterval:
			info.Interval = ad.Uint32()
		case tcaCodelECN:
			info.ECN = ad.Uint32()
		case tcaCodelCEThreshold:
			info.CEThreshold = ad.Uint32()
		default:
			return fmt.Errorf("extractCodelOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaFqUnspec = iota
	tcaFqPLimit
	tcaFqFlowPLimit
	tcaFqQuantum
	tcaFqInitQuantum
	tcaFqRateEnable
	tcaFqFlowDefaultRate
	tcaFqFlowMaxRate
	tcaFqBucketsLog
	tcaFqFlowRefillDelay
	tcaFqOrphanMask
	tcaFqLowRateThreshold
	tcaFqCEThreshold
)

// Fq contains attributes of the fq discipline
type Fq struct {
	PLimit           uint32
	FlowPLimit       uint32
	Quantum          uint32
	InitQuantum      uint32
	RateEnable       uint32
	FlowDefaultRate  uint32
	FlowMaxRate      uint32
	BucketsLog       uint32
	FlowRefillDelay  uint32
	OrphanMask       uint32
	LowRateThreshold uint32
	CEThreshold      uint32
}

func extractFqOptions(data []byte, info *Fq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaFqPLimit:
			info.PLimit = ad.Uint32()
		case tcaFqFlowPLimit:
			info.FlowPLimit = ad.Uint32()
		case tcaFqQuantum:
			info.Quantum = ad.Uint32()
		case tcaFqInitQuantum:
			info.Quantum = ad.Uint32()
		case tcaFqRateEnable:
			info.RateEnable = ad.Uint32()
		case tcaFqFlowDefaultRate:
			info.FlowDefaultRate = ad.Uint32()
		case tcaFqFlowMaxRate:
			info.FlowMaxRate = ad.Uint32()
		case tcaFqBucketsLog:
			info.BucketsLog = ad.Uint32()
		case tcaFqFlowRefillDelay:
			info.FlowRefillDelay = ad.Uint32()
		case tcaFqOrphanMask:
			info.OrphanMask = ad.Uint32()
		case tcaFqLowRateThreshold:
			info.LowRateThreshold = ad.Uint32()
		case tcaFqCEThreshold:
			info.CEThreshold = ad.Uint32()
		default:
			return fmt.Errorf("extractFqOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaPieUnspec = iota
	tcaPieTarget
	tcaPieLimit
	tcaPieTUpdate
	tcaPieAlpha
	tcaPieBeta
	tcaPieECN
	tcaPieBytemode
)

// Pie contains attributes of the pie discipline
type Pie struct {
	Target   uint32
	Limit    uint32
	TUpdate  uint32
	Alpha    uint32
	Beta     uint32
	ECN      uint32
	Bytemode uint32
}

func extractPieOptions(data []byte, info *Pie) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaPieTarget:
			info.Target = ad.Uint32()
		case tcaPieLimit:
			info.Limit = ad.Uint32()
		case tcaPieTUpdate:
			info.TUpdate = ad.Uint32()
		case tcaPieAlpha:
			info.Alpha = ad.Uint32()
		case tcaPieBeta:
			info.Beta = ad.Uint32()
		case tcaPieECN:
			info.ECN = ad.Uint32()
		case tcaPieBytemode:
			info.Bytemode = ad.Uint32()
		default:
			return fmt.Errorf("extractPieOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
