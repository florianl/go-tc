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

// GnetStatsBasic from include/uapi/linux/gen_stats.h
type GnetStatsBasic struct {
	Bytes   uint64
	Packets uint32
}

func extractGnetStatsBasic(data []byte, info *GnetStatsBasic) error {
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

// GnetStatsQueue from include/uapi/linux/gen_stats.h
type GnetStatsQueue struct {
	QueueLen   uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

func extractGnetStatsQueue(data []byte, info *GnetStatsQueue) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}
