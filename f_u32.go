package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaU32Unspec = iota
	tcaU32ClassID
	tcaU32Hash
	tcaU32Link
	tcaU32Divisor
	tcaU32Sel
	tcaU32Police
	tcaU32Act
	tcaU32InDev
	tcaU32Pcnt
	tcaU32Mark
	tcaU32Flags
	tcaU32Pad
)

// U32 contains attributes of the u32 discipline
type U32 struct {
	ClassID uint32
	Hash    uint32
	Link    uint32
	Divisor uint32
	Sel     *U32Sel
	InDev   string
	Pcnt    uint64
	Mark    *U32Mark
	Flags   uint32
	Police  *Police
}

// marshalU32 returns the binary encoding of U32
func marshalU32(info *U32) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("U32 options are missing")
	}

	// TODO: improve logic and check combinations

	if info.Sel != nil {
		data, err := validateU32SelOptions(info.Sel)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaU32Sel, Data: data})
	}

	if info.Mark != nil {
		data, err := validateU32MarkOptions(info.Mark)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaU32Mark, Data: data})
	}

	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaU32ClassID, Data: info.ClassID})
	}
	if info.Police != nil {
		data, err := marshalPolice(info.Police)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaU32Police, Data: data})
	}

	return marshalAttributes(options)
}

// unmarshalU32 parses the U32-encoded data and stores the result in the value pointed to by info.
func unmarshalU32(data []byte, info *U32) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaU32ClassID:
			info.ClassID = ad.Uint32()
		case tcaU32Hash:
			info.Hash = ad.Uint32()
		case tcaU32Link:
			info.Link = ad.Uint32()
		case tcaU32Divisor:
			info.Divisor = ad.Uint32()
		case tcaU32Sel:
			arg := &U32Sel{}
			if err := extractU32Sel(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Sel = arg
		case tcaU32Police:
			pol := &Police{}
			if err := unmarshalPolice(ad.Bytes(), pol); err != nil {
				return err
			}
			info.Police = pol
		case tcaU32InDev:
			info.InDev = ad.String()
		case tcaU32Pcnt:
			info.Pcnt = ad.Uint64()
		case tcaU32Mark:
			arg := &U32Mark{}
			if err := extractU32Mark(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Mark = arg
		case tcaU32Flags:
			info.Flags = ad.Uint32()
		case tcaU32Pad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalU32()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// U32Sel from include/uapi/linux/pkt_sched.h
type U32Sel struct {
	Flags    byte
	Offshift byte
	NKeys    byte
	OffMask  uint16
	Off      uint16
	Offoff   uint16
	Hoff     uint16
	Hmask    uint32
	U32Key
}

func validateU32SelOptions(info *U32Sel) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("U32Sel options are missing")
	}

	// TODO: improve logic and check combination
	return marshalAttributes(options)
}

func extractU32Sel(data []byte, info *U32Sel) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

//U32Mark from include/uapi/linux/pkt_sched.h
type U32Mark struct {
	Val     uint32
	Mask    uint32
	Success uint32
}

func validateU32MarkOptions(info *U32Mark) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}

func extractU32Mark(data []byte, info *U32Mark) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// U32Key from include/uapi/linux/pkt_sched.h
type U32Key struct {
	Mask    uint32
	Val     uint32
	Off     uint32
	OffMask uint32
}
