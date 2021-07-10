package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

type EmatchLayer uint8

const (
	EmatchLayerLink      = EmatchLayer(0)
	EmatchLayerNetwork   = EmatchLayer(1)
	EmatchLayerTransport = EmatchLayer(2)
)

type EmatchOpnd uint8

const (
	EmatchOpndEq = EmatchOpnd(0)
	EmatchOpndGt = EmatchOpnd(1)
	EmatchOpndLt = EmatchOpnd(2)
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
	Hdr        EmatchHdr
	U32Match   *U32Match
	CmpMatch   *CmpMatch
	IPSetMatch *IPSetMatch
}

// unmarshalEmatch parses the Ematch-encoded data and stores the result in the value pointed to by info.
func unmarshalEmatch(data []byte, info *Ematch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaEmatchTreeHdr:
			hdr := &EmatchTreeHdr{}
			err := unmarshalStruct(ad.Bytes(), hdr)
			concatError(multiError, err)
			info.Hdr = hdr
		case tcaEmatchTreeList:
			list := []EmatchMatch{}
			err := unmarshalEmatchTreeList(ad.Bytes(), &list)
			concatError(multiError, err)
			info.Matches = &list
		default:
			return fmt.Errorf("UnmarshalEmatch()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalEmatch returns the binary encoding of Ematch
func marshalEmatch(info *Ematch) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ematch: %w", ErrNoArg)
	}
	var multiError error

	if info.Hdr != nil {
		data, err := marshalStruct(info.Hdr)
		concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeHdr, Data: data})
	}
	if info.Matches != nil {
		data, err := marshalEmatchTreeList(info.Matches)
		concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEmatchTreeList | nlaFNnested, Data: data})
	}
	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}

func unmarshalEmatchTreeList(data []byte, info *[]EmatchMatch) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		match := EmatchMatch{}
		tmp := ad.Bytes()
		if err := unmarshalStruct(tmp[:8], &match.Hdr); err != nil {
			return err
		}
		switch match.Hdr.Kind {
		case EmatchU32:
			expr := &U32Match{}
			err := unmarshalU32Match(tmp[8:], expr)
			concatError(multiError, err)
			match.U32Match = expr
		case EmatchCmp:
			expr := &CmpMatch{}
			err := unmarshalCmpMatch(tmp[8:], expr)
			concatError(multiError, err)
			match.CmpMatch = expr
		case EmatchIPSet:
			expr := &IPSetMatch{}
			err := unmarshalIPSetMatch(tmp[8:], expr)
			concatError(multiError, err)
			match.IPSetMatch = expr
		default:
			return fmt.Errorf("unmarshalEmatchTreeList() kind %d is not yet implemented", match.Hdr.Kind)
		}
		*info = append(*info, match)
	}
	return concatError(multiError, ad.Err())
}

func marshalEmatchTreeList(info *[]EmatchMatch) ([]byte, error) {
	options := []tcOption{}

	for i, m := range *info {
		payload, err := marshalStruct(m.Hdr)
		if err != nil {
			return []byte{}, err
		}
		var expr []byte
		switch m.Hdr.Kind {
		case EmatchU32:
			expr, err = marshalU32Match(m.U32Match)
		case EmatchCmp:
			expr, err = marshalCmpMatch(m.CmpMatch)
		case EmatchIPSet:
			expr, err = marshalIPSetMatch(m.IPSetMatch)
		default:
			return []byte{}, fmt.Errorf("marshalEmatchTreeList() kind %d is not yet implemented", m.Hdr.Kind)
		}
		if err != nil {
			return []byte{}, fmt.Errorf("marshalEmatchTreeList(): %v", err)
		}
		payload = append(payload, expr...)
		options = append(options, tcOption{Interpretation: vtBytes, Type: uint16(i + 1), Data: payload})
	}
	return marshalAttributes(options)

}
