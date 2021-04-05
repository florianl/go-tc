package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaEmatchTreeUnspec = iota
	tcaEmatchTreeHdr
	tcaEmatchTreeList
)

type EmatchKind uint16

// Various Ematch kinds
const (
	EmatchContainer = EmatchKind(0)
	EmatchCmp       = EmatchKind(1)
	EmatchNByte     = EmatchKind(2)
	EmatchU32       = EmatchKind(3)
	EmatchMeta      = EmatchKind(4)
	EmatchText      = EmatchKind(5)
	EmatchVLan      = EmatchKind(6)
	EmatchCanID     = EmatchKind(7)
	EmatchIPSet     = EmatchKind(8)
	EmatchIPT       = EmatchKind(9)
)

// Ematch contains attributes of the ematch discipline
// https://man7.org/linux/man-pages/man8/tc-ematch.8.html
type Ematch struct {
	Hdr     *EmatchTreeHdr
	Matches *[]EmatchMatch
}

// tcf_ematch_tree_hdr from include/uapi/linux/pkt_cls.h
type EmatchTreeHdr struct {
	NMatches uint16
	ProgID   uint16
}

// tcf_ematch_hdr from include/uapi/linux/pkt_cls.h
type EmatchHdr struct {
	MatchID uint16
	Kind    EmatchKind
	Flags   uint16
	Pad     uint16
}

type EmatchMatch struct {
	Hdr  EmatchHdr
	Data []byte
}

// unmarshalEmatch parses the Ematch-encoded data and stores the result in the value pointed to by info.
func unmarshalEmatch(data []byte, info *Ematch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaEmatchTreeHdr:
			hdr := &EmatchTreeHdr{}
			if err := unmarshalStruct(ad.Bytes(), hdr); err != nil {
				return err
			}
			info.Hdr = hdr
		case tcaEmatchTreeList:
			list := []EmatchMatch{}
			if err := unmarshalEmatchTreeList(ad.Bytes(), &list); err != nil {
				return err
			}
			info.Matches = &list
		default:
			return fmt.Errorf("UnmarshalEmatch()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return ad.Err()
}

// marshalEmatch returns the binary encoding of Ematch
func marshalEmatch(info *Ematch) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ematch: %w", ErrNoArg)
	}

	if info.Hdr != nil {
		data, err := marshalStruct(info.Hdr)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeHdr, Data: data})
	}
	if info.Matches != nil {
		data, err := marshalEmatchTreeList(info.Matches)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeList | nlaFNnested, Data: data})
	}
	return marshalAttributes(options)
}

func unmarshalEmatchTreeList(data []byte, info *[]EmatchMatch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		match := EmatchMatch{}
		tmp := ad.Bytes()
		if err := unmarshalStruct(tmp[:8], &match.Hdr); err != nil {
			return err
		}
		match.Data = append(match.Data, tmp[8:]...)
		*info = append(*info, match)
	}
	return nil
}

func marshalEmatchTreeList(info *[]EmatchMatch) ([]byte, error) {
	options := []tcOption{}

	for i, m := range *info {
		payload, err := marshalStruct(m.Hdr)
		if err != nil {
			return []byte{}, err
		}
		payload = append(payload, m.Data...)
		options = append(options, tcOption{Interpretation: vtBytes, Type: uint16(i + 1), Data: payload})
	}
	return marshalAttributes(options)

}
