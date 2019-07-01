package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

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

// unmarshalFq parses the Fq-encoded data and stores the result in the value pointed to by info.
func unmarshalFq(data []byte, info *Fq) error {
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
			info.InitQuantum = ad.Uint32()
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
			return fmt.Errorf("unmarshalFq()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalFq returns the binary encoding of Fq
func marshalFq(info *Fq) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Fq options are missing")
	}

	// TODO: improve logic and check combinations
	if info.PLimit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqPLimit, Data: info.PLimit})
	}
	if info.FlowPLimit != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowPLimit, Data: info.FlowPLimit})
	}
	if info.Quantum != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqQuantum, Data: info.Quantum})
	}
	if info.InitQuantum != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqInitQuantum, Data: info.InitQuantum})
	}
	if info.RateEnable != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqRateEnable, Data: info.RateEnable})
	}
	if info.FlowDefaultRate != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowDefaultRate, Data: info.FlowDefaultRate})
	}
	if info.FlowMaxRate != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowMaxRate, Data: info.FlowMaxRate})
	}
	if info.BucketsLog != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqBucketsLog, Data: info.BucketsLog})
	}
	if info.FlowRefillDelay != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqFlowRefillDelay, Data: info.FlowRefillDelay})
	}
	if info.OrphanMask != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqOrphanMask, Data: info.OrphanMask})
	}
	if info.LowRateThreshold != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqLowRateThreshold, Data: info.LowRateThreshold})
	}
	if info.CEThreshold != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaFqCEThreshold, Data: info.CEThreshold})
	}
	return marshalAttributes(options)
}
