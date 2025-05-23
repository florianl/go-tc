package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaRedUnspec = iota
	tcaRedParms
	tcaRedStab
	tcaRedMaxP
	tcaRedFlags /* 32-bit bitfield value; 32-bit bitfield selector */
	tcaRedEarlyDropBlock
	tcaRedMarkBlock
)

// Red contains attributes of the red discipline
type Red struct {
	Parms          *RedQOpt
	MaxP           *uint32
	Flags          *uint64
	EarlyDropBlock *uint32
	MarkBlock      *uint32
}

// unmarshalRed parses the Red-encoded data and stores the result in the value pointed to by info.
func unmarshalRed(data []byte, info *Red) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaRedParms:
			opt := &RedQOpt{}
			multiError = unmarshalStruct(ad.Bytes(), opt)
			info.Parms = opt
		case tcaRedMaxP:
			info.MaxP = uint32Ptr(ad.Uint32())
		case tcaRedFlags:
			info.Flags = uint64Ptr(ad.Uint64())
		case tcaRedEarlyDropBlock:
			info.EarlyDropBlock = uint32Ptr(ad.Uint32())
		case tcaRedMarkBlock:
			info.MarkBlock = uint32Ptr(ad.Uint32())
		default:
			return fmt.Errorf("unmarshalRed()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalRed returns the binary encoding of Red
func marshalRed(info *Red) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Red: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRedParms, Data: data})
	}
	if info.MaxP != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRedMaxP, Data: uint32Value(info.MaxP)})
	}
	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaRedFlags, Data: uint64Value(info.Flags)})
	}
	if info.EarlyDropBlock != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRedEarlyDropBlock, Data: uint32Value(info.EarlyDropBlock)})
	}
	if info.MarkBlock != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRedMarkBlock, Data: uint32Value(info.MarkBlock)})
	}
	return marshalAttributes(options)
}

// RedQOpt from include/uapi/linux/pkt_sched.h
type RedQOpt struct {
	Limit    uint32
	QthMin   uint32
	QthMax   uint32
	Wlog     byte
	Plog     byte
	ScellLog byte
	Flags    byte
}
