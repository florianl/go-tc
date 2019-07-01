package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaAtmUnspec = iota
	tcaAtmFD
	tcaAtmPtr
	tcaAtmHdr
	tcaAtmExcess
	tcaAtmAddr
	tcaAtmState
)

// Atm contains attributes of the atm discipline
type Atm struct {
	FD     uint32
	Excess uint32
	Addr   *AtmPvc
	State  uint32
}

// unmarshalAtm parses the Atm-encoded data and stores the result in the value pointed to by info.
func unmarshalAtm(data []byte, info *Atm) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaAtmFD:
			info.FD = ad.Uint32()
		case tcaAtmExcess:
			info.Excess = ad.Uint32()
		case tcaAtmAddr:
			arg := &AtmPvc{}
			if err := extractAtmPvc(ad.Bytes(), arg); err != nil {
				return err
			}
			info.Addr = arg
		case tcaAtmState:
			info.State = ad.Uint32()
		default:
			return fmt.Errorf("unmarshalAtm()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return nil
}

// marshalAtm returns the binary encoding of Atm
func marshalAtm(info *Atm) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Atm options are missing")
	}
	// TODO: improve logic and check combinations

	if info.Addr != nil {
		data, err := validateAtmPvc(info.Addr)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaAtmAddr, Data: data})
	}
	if info.FD != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmFD, Data: info.FD})
	}
	if info.Excess != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmExcess, Data: info.Excess})
	}
	if info.State != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaAtmState, Data: info.State})
	}

	return marshalAttributes(options)
}

// AtmPvc from include/uapi/linux/atm.h
type AtmPvc struct {
	SapFamily byte
	Itf       byte
	Vpi       byte
	Vci       byte
}

func extractAtmPvc(data []byte, info *AtmPvc) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

func validateAtmPvc(info *AtmPvc) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}
