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
	case "qfq":
		info := &Qfq{}
		if err := extractQfqOptions(data, info); err != nil {
			return err
		}
		tc.Qfq = info
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

func extractClsact(data []byte) error {
	return fmt.Errorf("extractClsact()\t%v", data)
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
