package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

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

// ActBpf represents policing attributes of various filters and classes
type ActBpf struct {
	Tm     *Tcft
	Parms  *ActBpfParms
	Ops    []byte
	OpsLen uint16
	FD     uint32
	Name   string
	Tag    []byte
	ID     uint32
}

// unmarshalActBpf parses the ActBpf-encoded data and stores the result in the value pointed to by info.
func unmarshalActBpf(data []byte, info *ActBpf) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaActBpfTm:
			tm := &Tcft{}
			if err := unmarshalStruct(ad.Bytes(), tm); err != nil {
				return err
			}
			info.Tm = tm
		case tcaActBpfParms:
			parms := &ActBpfParms{}
			if err := unmarshalActBpfParms(ad.Bytes(), parms); err != nil {
				return err
			}
			info.Parms = parms
		case tcaActBpfOpsLen:
			info.OpsLen = ad.Uint16()
		case tcaActBpfOps:
			info.Ops = ad.Bytes()
		case tcaActBpfFD:
			info.FD = ad.Uint32()
		case tcaActBpfName:
			info.Name = ad.String()
		case tcaActBpfTag:
			info.Tag = ad.Bytes()
		case tcaActBpfPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalActBpf()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalActBpf returns the binary encoding of ActBpf
func marshalActBpf(info *ActBpf) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("ActBpf options are missing")
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if info.Name != "" {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaActBpfName, Data: info.Name})
	}
	if len(info.Tag) > 0 {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfTag, Data: info.Tag})
	}
	if info.FD != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActBpfFD, Data: info.FD})
	}
	if info.ID != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActBpfID, Data: info.ID})
	}
	if len(info.Ops) > 0 {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfOps, Data: info.Ops})
		options = append(options, tcOption{Interpretation: vtUint16, Type: tcaActBpfOpsLen, Data: info.OpsLen})
	}
	if info.Parms != nil {
		data, err := marshalActBpfParms(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActBpfParms, Data: data})
	}
	return marshalAttributes(options)
}

// ActBpfParms from include/uapi/linux/tc_act/tc_bpf.h
type ActBpfParms struct {
	Index   uint32
	Capab   uint32
	Action  uint32
	Refcnt  uint32
	Bindcnt uint32
}

// unmarshalActBpfParms parses the ActBpfParms-encoded data and stores the result in the value pointed to by info.
func unmarshalActBpfParms(data []byte, info *ActBpfParms) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

// marshalActBpfParms returns the binary encoding of ActBpfParms
func marshalActBpfParms(info *ActBpfParms) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}
