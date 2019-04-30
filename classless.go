package tc

import (
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
