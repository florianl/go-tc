package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
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

// BPF contains attributes of the bpf discipline
type BPF struct {
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      []byte
	ID       uint32
	Action   *Action
}

func validateBPFOptions(info *BPF) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("BPF options are missing")
	}

	// TODO: improve logic and check combinations

	if info.FD != 0 && len(info.Ops) != 0 {
		return []byte{}, fmt.Errorf("can not use FD and Ops at the same time")
	}

	if info.FD != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFd, Data: info.FD})
	}

	if len(info.Ops) != 0 {
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaBpfOpsLen, Data: info.OpsLen})
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfOps, Data: info.Ops})
	}

	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfClassid, Data: info.ClassID})
	}
	if len(info.Name) != 0 {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaBpfName, Data: info.Name})
	}
	if info.Flags != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlags, Data: info.Flags})
	}
	if info.FlagsGen != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlagsGen, Data: info.FlagsGen})
	}

	return marshalAttributes(options)

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
			info := &ActBpf{}
			if err := extractTcActBpf(ad.Bytes(), info); err != nil {
				return err
			}
			attr.Act = info
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

// BPFActionOptions contains various action attributes
type BPFActionOptions struct {
	OpsLen uint16
	Ops    []byte
	Tcft   *Tcft
	FD     uint32
	Name   string
	Act    *ActBpf
}

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

// ActionStats contains various statistics of a action
type ActionStats struct {
	Basic     *GenStatsBasic
	RateEst   *GenStatsRateEst
	Queue     *GenStatsQueue
	RateEst64 *GenStatsRateEst64
}

// Action describes a Traffic Control action
type Action struct {
	Kind       string
	Statistics *ActionStats
	BPFOptions *BPFActionOptions
}
