package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCakeUnspec = iota
	tcaCakePad
	tcaCakeBaseRate64
	tcaCakeDiffServMode
	tcaCakeAtm
	tcaCakeFlowMode
	tcaCakeOverhead
	tcaCakeRtt
	tcaCakeTarget
	tcaCakeAutorate
	tcaCakeMemory
	tcaCakeNat
	tcaCakeRaw
	tcaCakeWash
	tcaCakeMpu
	tcaCakeIngress
	tcaCakeAckFilter
	tcaCakeSplitGso
	tcaCakeFwMark
)

// Cake contains attributes of the cake discipline.
// http://man7.org/linux/man-pages/man8/tc-cake.8.html
type Cake struct {
	BaseRate     *uint64
	DiffServMode *uint32
	Atm          *uint32
	FlowMode     *uint32
	Overhead     *uint32
	Rtt          *uint32
	Target       *uint32
	Autorate     *uint32
	Memory       *uint32
	Nat          *uint32
	Raw          *uint32
	Wash         *uint32
	Mpu          *uint32
	Ingress      *uint32
	AckFilter    *uint32
	SplitGso     *uint32
	FwMark       *uint32
}

// unmarshalCake parses the Cake-encoded data and stores the result in the value pointed to by info.
func unmarshalCake(data []byte, info *Cake) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaCakeBaseRate64:
			tmp := ad.Uint64()
			info.BaseRate = &tmp
		case tcaCakeDiffServMode:
			tmp := ad.Uint32()
			info.DiffServMode = &tmp
		case tcaCakeAtm:
			tmp := ad.Uint32()
			info.Atm = &tmp
		case tcaCakeFlowMode:
			tmp := ad.Uint32()
			info.FlowMode = &tmp
		case tcaCakeOverhead:
			tmp := ad.Uint32()
			info.Overhead = &tmp
		case tcaCakeRtt:
			tmp := ad.Uint32()
			info.Rtt = &tmp
		case tcaCakeTarget:
			tmp := ad.Uint32()
			info.Target = &tmp
		case tcaCakeAutorate:
			tmp := ad.Uint32()
			info.Autorate = &tmp
		case tcaCakeMemory:
			tmp := ad.Uint32()
			info.Memory = &tmp
		case tcaCakeNat:
			tmp := ad.Uint32()
			info.Nat = &tmp
		case tcaCakeRaw:
			tmp := ad.Uint32()
			info.Raw = &tmp
		case tcaCakeWash:
			tmp := ad.Uint32()
			info.Wash = &tmp
		case tcaCakeMpu:
			tmp := ad.Uint32()
			info.Mpu = &tmp
		case tcaCakeIngress:
			tmp := ad.Uint32()
			info.Ingress = &tmp
		case tcaCakeAckFilter:
			tmp := ad.Uint32()
			info.AckFilter = &tmp
		case tcaCakeSplitGso:
			tmp := ad.Uint32()
			info.SplitGso = &tmp
		case tcaCakeFwMark:
			tmp := ad.Uint32()
			info.FwMark = &tmp
		case tcaCakePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalCake()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalCake returns the binary encoding of Red
func marshalCake(info *Cake) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Cake: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.BaseRate != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeBaseRate64, Data: *info.BaseRate})
	}
	if info.DiffServMode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeDiffServMode, Data: *info.DiffServMode})
	}
	if info.Atm != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAtm, Data: *info.Atm})
	}
	if info.FlowMode != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeFlowMode, Data: *info.FlowMode})
	}
	if info.Overhead != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeOverhead, Data: *info.Overhead})
	}
	if info.Rtt != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeRtt, Data: *info.Rtt})
	}
	if info.Target != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTarget, Data: *info.Target})
	}
	if info.Autorate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAutorate, Data: *info.Autorate})
	}
	if info.Memory != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeMemory, Data: *info.Memory})
	}
	if info.Nat != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeNat, Data: *info.Nat})
	}
	if info.Raw != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeRaw, Data: *info.Raw})
	}
	if info.Wash != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeWash, Data: *info.Wash})
	}
	if info.Mpu != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeMpu, Data: *info.Mpu})
	}
	if info.Ingress != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeIngress, Data: *info.Ingress})
	}
	if info.AckFilter != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeAckFilter, Data: *info.AckFilter})
	}
	if info.SplitGso != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeSplitGso, Data: *info.SplitGso})
	}
	if info.FwMark != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeFwMark, Data: *info.FwMark})
	}
	return marshalAttributes(options)
}

const (
	tcaCakeStatsInvalid = iota
	tcaCakeStatsPad
	tcaCakeStatsCapacityEstimate64
	tcaCakeStatsMemoryLimit
	tcaCakeStatsMemoryUsed
	tcaCakeStatsAvgNetoff
	tcaCakeStatsMinNetlen
	tcaCakeStatsMaxNetlen
	tcaCakeStatsMinAdjlen
	tcaCakeStatsMaxAdjlen
	tcaCakeStatsTinStats
	tcaCakeStatsDeficit
	tcaCakeStatsCobaltCount
	tcaCakeStatsDropping
	tcaCakeStatsDropNextUs
	tcaCakeStatsPDrop
	tcaCakeStatsBlueTimerUs
	tcaCakeStatsActiveQueues
)

const (
	tcaCakeTinStatsUnspec = iota
	tcaCakeTinStatsPad
	tcaCakeTinStatsSentPackets
	tcaCakeTinStatsSentBytes64
	tcaCakeTinStatsDroppedPackets
	tcaCakeTinStatsDroppedBytes64
	tcaCakeTinStatsAcksDroppedPackets
	tcaCakeTinStatsAcksDroppedBytes64
	tcaCakeTinStatsEcnMarkedPackets
	tcaCakeTinStatsEcnMarkedBytes64
	tcaCakeTinStatsBacklogPackets
	tcaCakeTinStatsBacklogBytes
	tcaCakeTinStatsThresholdRate64
	tcaCakeTinStatsTargetUs
	tcaCakeTinStatsIntervalUs
	tcaCakeTinStatsWayIndirectHits
	tcaCakeTinStatsWayMisses
	tcaCakeTinStatsWayCollisions
	tcaCakeTinStatsPeakDelayUs
	tcaCakeTinStatsAvgDelayUs
	tcaCakeTinStatsBaseDelayUs
	tcaCakeTinStatsSparseFlows
	tcaCakeTinStatsBulkFlows
	tcaCakeTinStatsUnresponsiveFlows
	tcaCakeTinStatsMaxSkblen
	tcaCakeTinStatsFlowQuantum
)

// CakeXStats contains the extended statistics of the cake discipline
// (TCA_CAKE_STATS_* from include/uapi/linux/pkt_sched.h).
// Qdisc dumps set the qdisc wide values and Tins, class dumps set the
// per flow values (Deficit, CobaltCount, Dropping, DropNextUs, PDrop and BlueTimerUs).
type CakeXStats struct {
	CapacityEstimate *uint64
	MemoryLimit      *uint32
	MemoryUsed       *uint32
	AvgNetoff        *uint32
	MinNetlen        *uint32
	MaxNetlen        *uint32
	MinAdjlen        *uint32
	MaxAdjlen        *uint32
	ActiveQueues     *uint32
	Tins             []CakeTinStats
	Deficit          *int32
	CobaltCount      *uint32
	Dropping         *uint32
	DropNextUs       *int32
	PDrop            *uint32
	BlueTimerUs      *int32
}

// CakeTinStats contains the statistics of a single cake tin
// (TCA_CAKE_TIN_STATS_* from include/uapi/linux/pkt_sched.h).
type CakeTinStats struct {
	SentPackets        *uint32
	SentBytes          *uint64
	DroppedPackets     *uint32
	DroppedBytes       *uint64
	AcksDroppedPackets *uint32
	AcksDroppedBytes   *uint64
	EcnMarkedPackets   *uint32
	EcnMarkedBytes     *uint64
	BacklogPackets     *uint32
	BacklogBytes       *uint32
	ThresholdRate      *uint64
	TargetUs           *uint32
	IntervalUs         *uint32
	WayIndirectHits    *uint32
	WayMisses          *uint32
	WayCollisions      *uint32
	PeakDelayUs        *uint32
	AvgDelayUs         *uint32
	BaseDelayUs        *uint32
	SparseFlows        *uint32
	BulkFlows          *uint32
	UnresponsiveFlows  *uint32
	MaxSkblen          *uint32
	FlowQuantum        *uint32
}

// unmarshalCakeXStats parses the TCA_STATS_APP payload of cake and stores the result in the value pointed to by info.
// Unknown attributes are skipped, so newer kernels exporting additional statistics do not break parsing.
func unmarshalCakeXStats(data []byte, info *CakeXStats) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaCakeStatsCapacityEstimate64:
			info.CapacityEstimate = uint64Ptr(ad.Uint64())
		case tcaCakeStatsMemoryLimit:
			info.MemoryLimit = uint32Ptr(ad.Uint32())
		case tcaCakeStatsMemoryUsed:
			info.MemoryUsed = uint32Ptr(ad.Uint32())
		case tcaCakeStatsAvgNetoff:
			info.AvgNetoff = uint32Ptr(ad.Uint32())
		case tcaCakeStatsMinNetlen:
			info.MinNetlen = uint32Ptr(ad.Uint32())
		case tcaCakeStatsMaxNetlen:
			info.MaxNetlen = uint32Ptr(ad.Uint32())
		case tcaCakeStatsMinAdjlen:
			info.MinAdjlen = uint32Ptr(ad.Uint32())
		case tcaCakeStatsMaxAdjlen:
			info.MaxAdjlen = uint32Ptr(ad.Uint32())
		case tcaCakeStatsActiveQueues:
			info.ActiveQueues = uint32Ptr(ad.Uint32())
		case tcaCakeStatsTinStats:
			err := unmarshalCakeTins(ad.Bytes(), info)
			multiError = concatError(multiError, err)
		case tcaCakeStatsDeficit:
			info.Deficit = int32Ptr(ad.Int32())
		case tcaCakeStatsCobaltCount:
			info.CobaltCount = uint32Ptr(ad.Uint32())
		case tcaCakeStatsDropping:
			info.Dropping = uint32Ptr(ad.Uint32())
		case tcaCakeStatsDropNextUs:
			info.DropNextUs = int32Ptr(ad.Int32())
		case tcaCakeStatsPDrop:
			info.PDrop = uint32Ptr(ad.Uint32())
		case tcaCakeStatsBlueTimerUs:
			info.BlueTimerUs = int32Ptr(ad.Int32())
		case tcaCakeStatsPad:
			// padding does not contain data, we just skip it
		default:
			// unknown statistics of newer kernels are skipped
		}
	}
	return concatError(multiError, ad.Err())
}

// unmarshalCakeTins parses the nested TCA_CAKE_STATS_TIN_STATS attribute.
// Each tin is a nested attribute of its own, whose type is the tin index + 1.
func unmarshalCakeTins(data []byte, info *CakeXStats) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		tin := CakeTinStats{}
		err := unmarshalCakeTinStats(ad.Bytes(), &tin)
		multiError = concatError(multiError, err)
		info.Tins = append(info.Tins, tin)
	}
	return concatError(multiError, ad.Err())
}

// unmarshalCakeTinStats parses the statistics of a single cake tin and stores the result in the value pointed to by info.
func unmarshalCakeTinStats(data []byte, info *CakeTinStats) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	for ad.Next() {
		switch ad.Type() {
		case tcaCakeTinStatsSentPackets:
			info.SentPackets = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsSentBytes64:
			info.SentBytes = uint64Ptr(ad.Uint64())
		case tcaCakeTinStatsDroppedPackets:
			info.DroppedPackets = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsDroppedBytes64:
			info.DroppedBytes = uint64Ptr(ad.Uint64())
		case tcaCakeTinStatsAcksDroppedPackets:
			info.AcksDroppedPackets = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsAcksDroppedBytes64:
			info.AcksDroppedBytes = uint64Ptr(ad.Uint64())
		case tcaCakeTinStatsEcnMarkedPackets:
			info.EcnMarkedPackets = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsEcnMarkedBytes64:
			info.EcnMarkedBytes = uint64Ptr(ad.Uint64())
		case tcaCakeTinStatsBacklogPackets:
			info.BacklogPackets = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsBacklogBytes:
			info.BacklogBytes = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsThresholdRate64:
			info.ThresholdRate = uint64Ptr(ad.Uint64())
		case tcaCakeTinStatsTargetUs:
			info.TargetUs = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsIntervalUs:
			info.IntervalUs = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsWayIndirectHits:
			info.WayIndirectHits = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsWayMisses:
			info.WayMisses = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsWayCollisions:
			info.WayCollisions = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsPeakDelayUs:
			info.PeakDelayUs = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsAvgDelayUs:
			info.AvgDelayUs = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsBaseDelayUs:
			info.BaseDelayUs = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsSparseFlows:
			info.SparseFlows = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsBulkFlows:
			info.BulkFlows = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsUnresponsiveFlows:
			info.UnresponsiveFlows = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsMaxSkblen:
			info.MaxSkblen = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsFlowQuantum:
			info.FlowQuantum = uint32Ptr(ad.Uint32())
		case tcaCakeTinStatsPad:
			// padding does not contain data, we just skip it
		default:
			// unknown statistics of newer kernels are skipped
		}
	}
	return ad.Err()
}

// marshalCakeXStats returns the binary encoding of CakeXStats
func marshalCakeXStats(info *CakeXStats) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("CakeXStats: %w", ErrNoArg)
	}

	if info.CapacityEstimate != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeStatsCapacityEstimate64, Data: *info.CapacityEstimate})
	}
	if info.MemoryLimit != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMemoryLimit, Data: *info.MemoryLimit})
	}
	if info.MemoryUsed != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMemoryUsed, Data: *info.MemoryUsed})
	}
	if info.AvgNetoff != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsAvgNetoff, Data: *info.AvgNetoff})
	}
	if info.MinNetlen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMinNetlen, Data: *info.MinNetlen})
	}
	if info.MaxNetlen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMaxNetlen, Data: *info.MaxNetlen})
	}
	if info.MinAdjlen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMinAdjlen, Data: *info.MinAdjlen})
	}
	if info.MaxAdjlen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsMaxAdjlen, Data: *info.MaxAdjlen})
	}
	if info.ActiveQueues != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsActiveQueues, Data: *info.ActiveQueues})
	}
	if len(info.Tins) > 0 {
		tins := []tcOption{}
		for i := range info.Tins {
			data, err := marshalCakeTinStats(&info.Tins[i])
			if err != nil {
				return []byte{}, err
			}
			tins = append(tins, tcOption{Interpretation: vtBytes, Type: uint16(i + 1), Data: data})
		}
		data, err := marshalAttributes(tins)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCakeStatsTinStats, Data: data})
	}
	if info.Deficit != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaCakeStatsDeficit, Data: *info.Deficit})
	}
	if info.CobaltCount != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsCobaltCount, Data: *info.CobaltCount})
	}
	if info.Dropping != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsDropping, Data: *info.Dropping})
	}
	if info.DropNextUs != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaCakeStatsDropNextUs, Data: *info.DropNextUs})
	}
	if info.PDrop != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeStatsPDrop, Data: *info.PDrop})
	}
	if info.BlueTimerUs != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaCakeStatsBlueTimerUs, Data: *info.BlueTimerUs})
	}
	return marshalAttributes(options)
}

// marshalCakeTinStats returns the binary encoding of CakeTinStats
func marshalCakeTinStats(info *CakeTinStats) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("CakeTinStats: %w", ErrNoArg)
	}

	if info.SentPackets != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsSentPackets, Data: *info.SentPackets})
	}
	if info.SentBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeTinStatsSentBytes64, Data: *info.SentBytes})
	}
	if info.DroppedPackets != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsDroppedPackets, Data: *info.DroppedPackets})
	}
	if info.DroppedBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeTinStatsDroppedBytes64, Data: *info.DroppedBytes})
	}
	if info.AcksDroppedPackets != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsAcksDroppedPackets, Data: *info.AcksDroppedPackets})
	}
	if info.AcksDroppedBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeTinStatsAcksDroppedBytes64, Data: *info.AcksDroppedBytes})
	}
	if info.EcnMarkedPackets != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsEcnMarkedPackets, Data: *info.EcnMarkedPackets})
	}
	if info.EcnMarkedBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeTinStatsEcnMarkedBytes64, Data: *info.EcnMarkedBytes})
	}
	if info.BacklogPackets != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsBacklogPackets, Data: *info.BacklogPackets})
	}
	if info.BacklogBytes != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsBacklogBytes, Data: *info.BacklogBytes})
	}
	if info.ThresholdRate != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaCakeTinStatsThresholdRate64, Data: *info.ThresholdRate})
	}
	if info.TargetUs != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsTargetUs, Data: *info.TargetUs})
	}
	if info.IntervalUs != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsIntervalUs, Data: *info.IntervalUs})
	}
	if info.WayIndirectHits != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsWayIndirectHits, Data: *info.WayIndirectHits})
	}
	if info.WayMisses != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsWayMisses, Data: *info.WayMisses})
	}
	if info.WayCollisions != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsWayCollisions, Data: *info.WayCollisions})
	}
	if info.PeakDelayUs != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsPeakDelayUs, Data: *info.PeakDelayUs})
	}
	if info.AvgDelayUs != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsAvgDelayUs, Data: *info.AvgDelayUs})
	}
	if info.BaseDelayUs != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsBaseDelayUs, Data: *info.BaseDelayUs})
	}
	if info.SparseFlows != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsSparseFlows, Data: *info.SparseFlows})
	}
	if info.BulkFlows != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsBulkFlows, Data: *info.BulkFlows})
	}
	if info.UnresponsiveFlows != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsUnresponsiveFlows, Data: *info.UnresponsiveFlows})
	}
	if info.MaxSkblen != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsMaxSkblen, Data: *info.MaxSkblen})
	}
	if info.FlowQuantum != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaCakeTinStatsFlowQuantum, Data: *info.FlowQuantum})
	}
	return marshalAttributes(options)
}
