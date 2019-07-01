package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaRedUnspec = iota
	tcaRedParms
	tcaRedStab
	tcaRedMaxP
)

// Red contains attributes of the red discipline
type Red struct {
	Parms *RedQOpt
	MaxP  uint32
}

// unmarshalRed parses the Red-encoded data and stores the result in the value pointed to by info.
func unmarshalRed(data []byte, info *Red) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaRedParms:
			opt := &RedQOpt{}
			if err := extractRedQOpt(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		case tcaRedMaxP:
			info.MaxP = ad.Uint32()
		default:
			return fmt.Errorf("unmarshalRed()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalRed returns the binary encoding of Red
func marshalRed(info *Red) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Red options are missing")
	}

	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := validateRedQopt(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRedParms, Data: data})
	}
	if info.MaxP != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRedMaxP, Data: info.MaxP})
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

func extractRedQOpt(data []byte, info *RedQOpt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

func validateRedQopt(info *RedQOpt) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}
