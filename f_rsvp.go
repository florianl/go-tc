package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaRsvpUnspec = iota
	tcaRsvpClassID
	tcaRsvpDst
	tcaRsvpSrc
	tcaRsvpPInfo
	tcaRsvpPolice
	tcaRsvpAct
)

// Rsvp contains attributes of the rsvp discipline
type Rsvp struct {
	ClassID uint32
	Dst     []byte
	Src     []byte
	PInfo   *RsvpPInfo
	Police  *Police
}

//UnmarshalRsvp parses the Rsvp-encoded data and stores the result in the value pointed to by info.
func UnmarshalRsvp(data []byte, info *Rsvp) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaRsvpClassID:
			info.ClassID = ad.Uint32()
		case tcaRsvpDst:
			info.Dst = ad.Bytes()
		case tcaRsvpSrc:
			info.Src = ad.Bytes()
		case tcaRsvpPInfo:
			arg := &RsvpPInfo{}
			if err := extractRsvpPInfo(ad.Bytes(), arg); err != nil {
				return err
			}
			info.PInfo = arg
		case tcaRsvpPolice:
			pol := &Police{}
			if err := UnmarshalPolice(ad.Bytes(), pol); err != nil {
				return err
			}
			info.Police = pol
		default:
			return fmt.Errorf("UnmarshalRsvp()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil

}

// MarshalRsvp returns the binary encoding of Rsvp
func MarshalRsvp(info *Rsvp) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Rsvp options are missing")
	}

	// TODO: improve logic and check combinations
	if info.ClassID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4ClassID, Data: info.ClassID})
	}
	if info.PInfo != nil {
		data, err := validateRsvpPInfo(info.PInfo)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpPInfo, Data: data})
	}
	if info.Police != nil {
		data, err := MarshalPolice(info.Police)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpPolice, Data: data})
	}
	return marshalAttributes(options)
}

// RsvpPInfo from include/uapi/linux/pkt_sched.h
type RsvpPInfo struct {
	Dpi       RsvpGpi
	Spi       RsvpGpi
	Protocol  uint8
	TunnelID  uint8
	TunnelHdr uint8
	Pad       uint8
}

func extractRsvpPInfo(data []byte, info *RsvpPInfo) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

func validateRsvpPInfo(info *RsvpPInfo) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}

// RsvpGpi from include/uapi/linux/pkt_sched.h
type RsvpGpi struct {
	Key    uint32
	Mask   uint32
	Offset uint32
}
