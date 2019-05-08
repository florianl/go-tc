package tc

import (
	"bytes"
	"encoding/binary"
)

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

func extractTCStats(data []byte, info *Stats) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
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

func extractTCStats2(data []byte, info *Stats2) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// Tcft from include/uapi/linux/pkt_sched.h
type Tcft struct {
	Install  uint64
	LastUse  uint64
	Expires  uint64
	FirstUse uint64
}

func extractTcft(data []byte, info *Tcft) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// GenStatsBasic from include/uapi/linux/gen_stats.h
type GenStatsBasic struct {
	Bytes   uint64
	Packets uint32
}

func extractGnetStatsBasic(data []byte, info *GenStatsBasic) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// GenStatsRateEst from include/uapi/linux/gen_stats.h
type GenStatsRateEst struct {
	BytePerSecond   uint32
	PacketPerSecond uint32
}

func extractGenStatsRateEst(data []byte, info *GenStatsRateEst) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// GenStatsRateEst64 from include/uapi/linux/gen_stats.h
type GenStatsRateEst64 struct {
	BytePerSecond   uint64
	PacketPerSecond uint64
}

func extractGenStatsRateEst64(data []byte, info *GenStatsRateEst64) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// GenStatsQueue from include/uapi/linux/gen_stats.h
type GenStatsQueue struct {
	QueueLen   uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

func extractGnetStatsQueue(data []byte, info *GenStatsQueue) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// ActBpf from include/uapi/linux/tc_act/tc_bpf.h
type ActBpf struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	Refcnt  uint32
	Bindcnt uint32
}

func extractTcActBpf(data []byte, info *ActBpf) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
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

func extractRateSpec(data []byte, info *RateSpec) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// FifoOpt from from include/uapi/linux/pkt_sched.h
type FifoOpt struct {
	Limit uint32
}

func extractFifoOpt(data []byte, info *FifoOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}
