package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func unmarshalStruct(data []byte, s interface{}) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, s)
}

func marshalStruct(s interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, s)
	return buf.Bytes(), err
}

// Stats from include/uapi/linux/pkt_sched.h
type Stats struct {
	Bytes      uint64 /* Number of enqueued bytes */
	Packets    uint32 /* Number of enqueued packets	*/
	Drops      uint32 /* Packets dropped because of lack of resources */
	Overlimits uint32 /* Number of throttle events when this
	 * flow goes out of allocated bandwidth */
	Bps     uint32 /* Current flow byte rate */
	Pps     uint32 /* Current flow packet rate */
	Qlen    uint32
	Backlog uint32
}

// Stats2 from include/uapi/linux/pkt_sched.h
type Stats2 struct {
	// gnet_stats_basic
	Bytes   uint64
	Packets uint32
	//gnet_stats_queue
	Qlen       uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

// Tcft from include/uapi/linux/pkt_sched.h
type Tcft struct {
	Install  uint64
	LastUse  uint64
	Expires  uint64
	FirstUse uint64
}

// GenStatsBasic from include/uapi/linux/gen_stats.h
type GenStatsBasic struct {
	Bytes   uint64
	Packets uint32
}

// GenStatsRateEst from include/uapi/linux/gen_stats.h
type GenStatsRateEst struct {
	BytePerSecond   uint32
	PacketPerSecond uint32
}

// GenStatsRateEst64 from include/uapi/linux/gen_stats.h
type GenStatsRateEst64 struct {
	BytePerSecond   uint64
	PacketPerSecond uint64
}

// GenStatsQueue from include/uapi/linux/gen_stats.h
type GenStatsQueue struct {
	QueueLen   uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

// RateSpec from from include/uapi/linux/pkt_sched.h
type RateSpec struct {
	CellLog   uint8
	Linklayer uint8
	Overhead  uint16
	CellAlign uint16
	Mpu       uint16
	Rate      uint32
}

// Policy from from include/uapi/linux/pkt_sched.h
type Policy struct {
	Index    uint32
	Action   uint32
	Limit    uint32
	Burst    uint32
	Mtu      uint32
	Rate     RateSpec
	PeakRate RateSpec
	RefCnt   uint32
	BindCnt  uint32
	Capab    uint32
}

// FifoOpt from from include/uapi/linux/pkt_sched.h
type FifoOpt struct {
	Limit uint32
}

// SfqXStats from from include/uapi/linux/pkt_sched.h
type SfqXStats struct {
	Allot int32
}

// RedXStats from from include/uapi/linux/pkt_sched.h
type RedXStats struct {
	Early  uint32
	PDrop  uint32
	Other  uint32
	Marked uint32
}

// ChokeXStats from from include/uapi/linux/pkt_sched.h
type ChokeXStats struct {
	Early   uint32
	PDrop   uint32
	Other   uint32
	Marked  uint32
	Matched uint32
}

// HtbXStats from from include/uapi/linux/pkt_sched.h
type HtbXStats struct {
	Lends   uint32
	Borrows uint32
	Giants  uint32
	Tokens  uint32
	CTokens uint32
}

// CbqXStats from from include/uapi/linux/pkt_sched.h
type CbqXStats struct {
	Borrows     uint32
	Overactions uint32
	AvgIdle     int32
	Undertime   int32
}

// SfbXStats from from include/uapi/linux/pkt_sched.h
type SfbXStats struct {
	EarlyDrop   uint32
	PenaltyDrop uint32
	BucketDrop  uint32
	QueueDrop   uint32
	ChildDrop   uint32
	Marked      uint32
	MaxQlen     uint32
	MaxProb     uint32
	AvgProb     uint32
}

// CodelXStats from from include/uapi/linux/pkt_sched.h
type CodelXStats struct {
	MaxPacket     uint32
	Count         uint32
	LastCount     uint32
	LDelay        uint32
	DropNext      int32
	DropOverlimit uint32
	EcnMark       uint32
	Dropping      uint32
	CeMark        uint32
}

// HhfXStats from from include/uapi/linux/pkt_sched.h
type HhfXStats struct {
	DropOverlimit uint32
	HhOverlimit   uint32
	HhTotCount    uint32
	HhCurCount    uint32
}

// PieXStats from from include/uapi/linux/pkt_sched.h
type PieXStats struct {
	Prob      uint64
	Delay     uint32
	AvgDqRate uint32
	PacketsIn uint32
	Dropped   uint32
	Overlimit uint32
	Maxq      uint32
	EcnMark   uint32
}

// FqCodelQdStats from from include/uapi/linux/pkt_sched.h
type FqCodelQdStats struct {
	MaxPacket      uint32
	DropOverlimit  uint32
	EcnMark        uint32
	NewFlowCount   uint32
	NewFlowsLen    uint32
	OldFlowsLen    uint32
	CeMark         uint32
	MemoryUsage    uint32
	DropOvermemory uint32
}

// FqCodelClStats from from include/uapi/linux/pkt_sched.h
type FqCodelClStats struct {
	Deficit   int32
	LDelay    uint32
	Count     uint32
	LastCount uint32
	Dropping  uint32
	DropNext  int32
}

// FqCodelXStats from from include/uapi/linux/pkt_sched.h
type FqCodelXStats struct {
	Type uint32
	Qd   *FqCodelQdStats
	Cl   *FqCodelClStats
}

func extractFqCodelXStats(data []byte, info *FqCodelXStats) error {
	info.Type = nativeEndian.Uint32(data[:4])
	switch info.Type {
	case tcaFqCodelXStatsQdisc:
		b := bytes.NewReader(data[4:])
		stats := &FqCodelQdStats{}
		if err := binary.Read(b, nativeEndian, stats); err != nil {
			return err
		}
		info.Qd = stats
	case tcaFqCodelXStatsClass:
		b := bytes.NewReader(data[4:])
		stats := &FqCodelClStats{}
		if err := binary.Read(b, nativeEndian, stats); err != nil {
			return err
		}
		info.Cl = stats
	default:
		return fmt.Errorf("extractFqCodelXStats(): unsupported type: %d", info.Type)
	}
	return nil
}
func marshalFqCodelXStats(v *FqCodelXStats) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, v.Type)
	if err != nil {
		return []byte{}, err
	}
	var subStat []byte
	switch v.Type {
	case tcaFqCodelXStatsQdisc:
		subStat, err = marshalStruct(v.Qd)
	case tcaFqCodelXStatsClass:
		subStat, err = marshalStruct(v.Cl)
	default:
		err = fmt.Errorf("marshalFqCodelXStats(): unknown FqCodelXStat type: %d", v.Type)
	}
	if err != nil {
		return []byte{}, err
	}
	_, err = buf.Write(subStat)
	return buf.Bytes(), err
}
