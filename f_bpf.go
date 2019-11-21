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

// Bpf contains attributes of the bpf discipline
type Bpf struct {
	Action   *Action
	Police   *Police
	ClassID  uint32
	OpsLen   uint16
	Ops      []byte
	FD       uint32
	Name     string
	Flags    uint32
	FlagsGen uint32
	Tag      []byte
	ID       uint32
}

// unmarshalBpf parses the Bpf-encoded data and stores the result in the value pointed to by info.
func unmarshalBpf(data []byte, info *Bpf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaBpfPolice:
			pol := &Police{}
			if err := unmarshalPolice(ad.Bytes(), pol); err != nil {
				return err
			}
			info.Police = pol
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
			return fmt.Errorf("unmarshalBpf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalBpf returns the binary encoding of Bpf
func marshalBpf(info *Bpf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Bpf: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if len(info.Ops) > 0 {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfOps, Data: info.Ops})
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaBpfOpsLen, Data: info.OpsLen})
	}
	if info.Name != "" {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaBpfName, Data: info.Name})
	}
	if info.FD != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFd, Data: info.FD})
	}
	if info.ID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfID, Data: info.ID})
	}
	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfClassid, Data: info.ClassID})
	}
	if len(info.Tag) > 0 {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaBpfTag, Data: info.Tag})
	}
	if info.Flags != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlags, Data: info.Flags})
	}
	if info.FlagsGen != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaBpfFlagsGen, Data: info.FlagsGen})
	}
	return marshalAttributes(options)
}
