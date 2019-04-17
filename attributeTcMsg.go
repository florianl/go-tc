package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

func extractTcmsgAttributes(data []byte, info *Attribute) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var kind string
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaKind:
			info.Kind = ad.String()
			kind = ad.String()
		case tcaOptions:
			if err := extractTCAOptions(ad.Bytes(), info, kind); err != nil {
				return err
			}
		case tcaChain:
			info.Chain = ad.Uint32()
		case tcaXstats:
			tcstats := &Stats{}
			if err := extractTCStats(ad.Bytes(), tcstats); err != nil {
				return err
			}
			info.XStats = tcstats
		case tcaStats:
			tcstats := &Stats{}
			if err := extractTCStats(ad.Bytes(), tcstats); err != nil {
				return err
			}
			info.Stats = tcstats
		case tcaStats2:
			tcstats2 := &Stats2{}
			if err := extractTCStats2(ad.Bytes(), tcstats2); err != nil {
				return err
			}
			info.Stats2 = tcstats2
		case tcaHwOffload:
			info.HwOffload = ad.Uint8()
		case tcaEgressBlock:
			info.EgressBlock = ad.Uint32()
		case tcaIngressBlock:
			info.IngressBlock = ad.Uint32()
		default:
			return fmt.Errorf("extractTcmsgAttributes()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return nil
}

func extractTCAOptions(data []byte, tc *Attribute, kind string) error {
	switch kind {
	case "fq_codel":
		info := &FqCodel{}
		if err := extractFqCodelOptions(data, info); err != nil {
			return err
		}
		tc.FqCodel = info
	case "clsact":
		return extractClsact(data)
	case "bpf":
		info := &BPF{}
		if err := extractBpfOptions(data, info); err != nil {
			return err
		}
		tc.BPF = info
	default:
		return fmt.Errorf("extractTCAOptions(): unsupported kind: %s", kind)
	}

	return nil
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

func extractTcAction(data []byte, act *Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaActKind:
			act.Kind = ad.String()
		case tcaActOptions:
			options := &BPFActionOptions{}
			if err := extractActBpfOptions(ad.Bytes(), options); err != nil {
				return err
			}
			act.BPFOptions = options
		case tcaActStats:
			stats := &ActionStats{}
			if err := extractActStats(ad.Bytes(), stats); err != nil {
				return err
			}
			act.Statistics = stats
		default:
			return fmt.Errorf("extractTcAction()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}

	return nil
}

func extractActStats(data []byte, stats *ActionStats) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaStatsBasic:
			info := &GenStatsBasic{}
			if err := extractGnetStatsBasic(ad.Bytes(), info); err != nil {
				return err
			}
			stats.Basic = info
		case tcaStatsRateEst:
			info := &GenStatsRateEst{}
			if err := extractGenStatsRateEst(ad.Bytes(), info); err != nil {
				return err
			}
			stats.RateEst = info
		case tcaStatsQueue:
			info := &GenStatsQueue{}
			if err := extractGnetStatsQueue(ad.Bytes(), info); err != nil {
				return err
			}
			stats.Queue = info
		case tcaStatsRateEst64:
			info := &GenStatsRateEst64{}
			if err := extractGenStatsRateEst64(ad.Bytes(), info); err != nil {
				return err
			}
			stats.RateEst64 = info
		case tcaStatsBasicHw:
			// ignore it for the moment. TODO
		default:
			return fmt.Errorf("extractActStats()\t%d\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
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
			/* struct tc_act_bpf */
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

func extractClsact(data []byte) error {
	return fmt.Errorf("extractClsact()\t%v", data)
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
	tcaUnspec = iota
	tcaKind
	tcaOptions
	tcaStats
	tcaXstats
	tcaRate
	tcaFcnt
	tcaStats2
	tcaStab
	tcaPad
	tcaDumpInvisible
	tcaChain
	tcaHwOffload
	tcaIngressBlock
	tcaEgressBlock
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

const (
	tcaActUnspec = iota
	tcaActKind
	tcaActOptions
	tcaActIndex
	tcaActStats
	tcaActPad
	tcaActCookie
)

const (
	tcaStatsUnspec = iota
	tcaStatsBasic
	tcaStatsRateEst
	tcaStatsQueue
	tcaStatsApp
	tcaStatsRateEst64
	tcaStatsPAD
	tcaStatsBasicHw
)

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
