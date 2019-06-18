//+build linux

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
	var options []byte
	var xStats []byte
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaKind:
			info.Kind = ad.String()
		case tcaOptions:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			options = ad.Bytes()
		case tcaChain:
			info.Chain = ad.Uint32()
		case tcaXstats:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			xStats = ad.Bytes()
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
	if len(options) > 0 {
		if err := extractTCAOptions(options, info, info.Kind); err != nil {
			return err
		}
	}
	if len(xStats) > 0 {
		tcxstats := &XStats{}
		if err := extractXStats(ad.Bytes(), tcxstats, info.Kind); err != nil {
			return err
		}
		info.XStats = tcxstats
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
	case "codel":
		info := &Codel{}
		if err := extractCodelOptions(data, info); err != nil {
			return err
		}
		tc.Codel = info
	case "fq":
		info := &Fq{}
		if err := extractFqOptions(data, info); err != nil {
			return err
		}
		tc.Fq = info
	case "pie":
		info := &Pie{}
		if err := extractPieOptions(data, info); err != nil {
			return err
		}
		tc.Pie = info
	case "hhf":
		info := &Hhf{}
		if err := extractHhfOptions(data, info); err != nil {
			return err
		}
		tc.Hhf = info
	case "htb":
		info := &Htb{}
		if err := extractHtbOptions(data, info); err != nil {
			return err
		}
		tc.Htb = info
	case "hfsc":
		info := &Hfsc{}
		if err := extractHfscOptions(data, info); err != nil {
			return err
		}
		tc.Hfsc = info
	case "dsmark":
		info := &Dsmark{}
		if err := extractDsmarkOptions(data, info); err != nil {
			return err
		}
		tc.Dsmark = info
	case "drr":
		info := &Drr{}
		if err := extractDrrOptions(data, info); err != nil {
			return err
		}
		tc.Drr = info
	case "cbq":
		info := &Cbq{}
		if err := extractCbqOptions(data, info); err != nil {
			return err
		}
		tc.Cbq = info
	case "atm":
		info := &Atm{}
		if err := extractAtmOptions(data, info); err != nil {
			return err
		}
		tc.Atm = info
	case "tbf":
		info := &Tbf{}
		if err := extractTbfOptions(data, info); err != nil {
			return err
		}
		tc.Tbf = info
	case "sfb":
		info := &Sfb{}
		if err := extractSfbOptions(data, info); err != nil {
			return err
		}
		tc.Sfb = info
	case "red":
		info := &Red{}
		if err := extractRedOptions(data, info); err != nil {
			return err
		}
		tc.Red = info
	case "pfifo":
		limit := &FifoOpt{}
		if err := extractFifoOpt(data, limit); err != nil {
			return err
		}
		tc.Pfifo = limit
	case "mqprio":
		info := &MqPrio{}
		if err := extractMqPrioOptions(data, info); err != nil {
			return err
		}
		tc.MqPrio = info
	case "bfifo":
		limit := &FifoOpt{}
		if err := extractFifoOpt(data, limit); err != nil {
			return err
		}
		tc.Bfifo = limit
	case "clsact":
		return extractClsact(data)
	case "ingress":
		return extractIngress(data)
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
	case "u32":
		info := &U32{}
		if err := UnmarshalU32(data, info); err != nil {
			return err
		}
		tc.U32 = info
	case "rsvp":
		info := &Rsvp{}
		if err := extractRsvpOptions(data, info); err != nil {
			return err
		}
		tc.Rsvp = info
	case "route":
		info := &Route4{}
		if err := extractRoute4Options(data, info); err != nil {
			return err
		}
		tc.Route4 = info
	case "fw":
		info := &Fw{}
		if err := extractFwOptions(data, info); err != nil {
			return err
		}
		tc.Fw = info
	case "flow":
		info := &Flow{}
		if err := extractFlowOptions(data, info); err != nil {
			return err
		}
		tc.Flow = info
	default:
		return fmt.Errorf("extractTCAOptions(): unsupported kind: %s", kind)
	}

	return nil
}

func extractXStats(data []byte, tc *XStats, kind string) error {
	switch kind {
	case "sfq":
		info := &SfqXStats{}
		if err := extractSfqXStats(data, info); err != nil {
			return err
		}
		tc.Sfq = info
	case "sfb":
		info := &SfbXStats{}
		if err := extractSfbXStats(data, info); err != nil {
			return err
		}
		tc.Sfb = info
	case "red":
		info := &RedXStats{}
		if err := extractRedXStats(data, info); err != nil {
			return err
		}
		tc.Red = info
	case "choke":
		info := &ChokeXStats{}
		if err := extractChokeXStats(data, info); err != nil {
			return err
		}
		tc.Choke = info
	case "htb":
		info := &HtbXStats{}
		if err := extractHtbXStats(data, info); err != nil {
			return err
		}
		tc.Htb = info
	case "cbq":
		info := &CbqXStats{}
		if err := extractCbqXStats(data, info); err != nil {
			return err
		}
		tc.Cbq = info
	case "codel":
		info := &CodelXStats{}
		if err := extractCodelXStats(data, info); err != nil {
			return err
		}
		tc.Codel = info
	case "hhf":
		info := &HhfXStats{}
		if err := extractHhfXStats(data, info); err != nil {
			return err
		}
		tc.Hhf = info
	case "pie":
		info := &PieXStats{}
		if err := extractPieXStats(data, info); err != nil {
			return err
		}
		tc.Pie = info
	case "fq_codel":
		info := &FqCodelXStats{}
		if err := extractFqCodelXStats(data, info); err != nil {
			return err
		}
		tc.FqCodel = info
	default:
		return fmt.Errorf("extractXStats(): unsupported kind: %s", kind)
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
	// Clsact is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractClsact()\t%v", data)
	}
	return nil
}

func extractIngress(data []byte) error {
	// Ingress is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractIngress()\t%v", data)
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
