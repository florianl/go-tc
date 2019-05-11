//+build linux

package tc

import (
	"github.com/mdlayher/netlink"
)

type valueType int

const (
	vtUint8 valueType = iota
	vtUint16
	vtUint32
	vtUint64
	vtString
	vtBytes
)

type tcOption struct {
	Interpretation valueType
	Type           uint16
	Data           interface{}
}

func marshalAttributes(options []tcOption) ([]byte, error) {
	ad := netlink.NewAttributeEncoder()

	for _, option := range options {
		switch option.Interpretation {
		case vtUint8:
			ad.Uint8(option.Type, (option.Data).(uint8))
		case vtUint16:
			ad.Uint16(option.Type, (option.Data).(uint16))
		case vtUint32:
			ad.Uint32(option.Type, (option.Data).(uint32))
		case vtUint64:
			ad.Uint64(option.Type, (option.Data).(uint64))
		case vtString:
			ad.String(option.Type, (option.Data).(string))
		case vtBytes:
			ad.Bytes(option.Type, (option.Data).([]byte))
		}
	}

	return ad.Encode()
}
